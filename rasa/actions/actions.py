# This files contains your custom actions which can be used to run
# custom Python code.
#
# See this guide on how to implement these action:
# https://rasa.com/docs/rasa/custom-actions

from rasa_sdk import Action, Tracker
from rasa_sdk.executor import CollectingDispatcher
from typing import Dict, Text, Any, List, Union, Optional
from influxdb import InfluxDBClient
from dateutil import relativedelta, parser


client = InfluxDBClient('localhost', 8086, '', '', 'rzpftx')

measurements = {
    "authorized_payments": "payment_authorized",
    "failed_payments": "payment_failed"
}


class ActionListMetrics(Action):
    def name(self) -> Text:
        return "action_list_metrics"

    def run(self,
            dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        result = ["authorized_payments (count/volume)",
                  "failed_payments (count/volume)", "success_rate"]

        dispatcher.utter_message(
            "The following metrics are available to query from: {}".format(", ".join(result)))
        return []


class ActionListDimensions(Action):

    def name(self) -> Text:
        return "action_list_dimensions"

    def run(self,
            dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        measurement = next(
            tracker.get_latest_entity_values("measurement"), None)
        if measurement not in ["authorized_payments", "failed_payments"]:
            dispatcher.utter_message(
                "Could not recognize metric: {}".format(measurement))
            return []

        result = client.query(
            "SHOW TAG KEYS FROM {}".format(measurements[measurement]))
        points = list(result.get_points())
        tags = map(lambda x: x["tagKey"], points)
        dispatcher.utter_message("The following dimensions are available for {}: {}".format(
            measurement, ", ".join(tags)))
        return []


class ActionCountMeasurements(Action):
    def name(self) -> Text:
        return "action_count_measurements"

    def run(self,
            dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        measurement = next(
            tracker.get_latest_entity_values("measurement"), None)
        time = self.parse_time(tracker.latest_message)

        if measurement not in ["authorized_payments", "failed_payments"]:
            dispatcher.utter_message(
                "Could not recognize metric: {}".format(measurement))
            return []

        query = "SELECT COUNT(*) FROM {} where time >= '{}' and time <= '{}'".format(
            measurements[measurement], time["start_time"], time["end_time"])
        print(query)

        result = client.query(query)
        points = list(result.get_points())
        count = points[0].get("count_amount", 0) if (len(points) == 1) else 0
        dispatcher.utter_message("Count: {}".format(count))
        return []

    def parse_time(self, latest_message):
        timeinfo = next((x for x in latest_message.get(
            "entities") if x["entity"] == "time"), {}).get("additional_info", None)
        if timeinfo.get("type") == "intervalisoformat":
            return self.close_interval_duckling_time(timeinfo)
        elif timeinfo.get("type") == "value":
            return self.make_interval_from_value_duckling_time(timeinfo)

    def close_interval_duckling_time(self,
                                     timeinfo: Dict[Text, Any]
                                     ) -> Optional[Dict[Text, Any]]:
        grain = timeinfo.get("to", timeinfo.get("from", {})).get("grain")
        start = timeinfo.get("from", {}).get("value")
        end = timeinfo.get("to", {}).get("value")
        if (start or end) and not (start and end):
            deltaargs = {f"{grain}s": 1}
            delta = relativedelta.relativedelta(**deltaargs)
            if start:
                parsedstart = parser.isoparse(start)
                parsedend = parsedstart + delta
                end = parsedend.isoformat()
            elif end:
                parsedend = parser.isoparse(end)
                parsedstart = parsedend - delta
                start = parsedstart.isoformat()
        return {
            "start_time": start,
            "end_time": end,
            "grain": grain
        }

    def make_interval_from_value_duckling_time(
        self, timeinfo: Dict[Text, Any]
    ) -> Dict[Text, Any]:
        grain = timeinfo.get("grain")
        start = timeinfo.get("value")
        parsedstart = parser.isoparse(start)
        deltaargs = {f"{grain}s": 1}
        delta = relativedelta.relativedelta(**deltaargs)
        parsedend = parsedstart + delta
        end = parsedend.isoformat()
        return {
            "start_time": start,
            "end_time": end,
            "grain": grain,
        }
