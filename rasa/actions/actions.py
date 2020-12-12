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
from actions.parser import parse_time


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

        method = next(
            tracker.get_latest_entity_values("method"), None)

        time = parse_time(tracker.latest_message)

        if measurement not in ["authorized_payments", "failed_payments"]:
            dispatcher.utter_message(
                "Could not recognize metric: {}".format(measurement))
            return []
        filter = "and method = '{}'".format(
            method) if (method is not None) else ""

        query = "SELECT COUNT(*) FROM {} where time >= '{}' and time <= '{}' {}".format(
            measurements[measurement], time["start_time"], time["end_time"], filter)

        print(query)

        result = client.query(query)
        points = list(result.get_points())
        count = points[0].get("count_amount", 0) if (len(points) == 1) else 0
        dispatcher.utter_message("Count: {}".format(count))
        return []


class ActionSumMeasurements(Action):
    def name(self) -> Text:
        return "action_sum_measurements"

    def run(self,
            dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        measurement = next(
            tracker.get_latest_entity_values("measurement"), None)

        method = next(
            tracker.get_latest_entity_values("method"), None)

        filter = "and method = '{}'".format(
            method) if (method is not None) else ""

        time = parse_time(tracker.latest_message)

        if measurement not in ["authorized_payments", "failed_payments"]:
            dispatcher.utter_message(
                "Could not recognize metric: {}".format(measurement))
            return []

        query = "SELECT SUM(amount) as amount FROM {} where time >= '{}' and time <= '{}' {}".format(
            measurements[measurement], time["start_time"], time["end_time"], filter)
        
        print(query)

        result = client.query(query)
        points = list(result.get_points())
        sum = points[0].get("amount", 0) if (len(points) == 1) else 0
        dispatcher.utter_message("INR: {}".format(sum/100))
        return []
