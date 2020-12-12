# This files contains your custom actions which can be used to run
# custom Python code.
#
# See this guide on how to implement these action:
# https://rasa.com/docs/rasa/custom-actions

from rasa_sdk import Action, Tracker
from rasa_sdk.executor import CollectingDispatcher
from rasa_sdk.events import SlotSet
from typing import Dict, Text, Any, List, Union, Optional
from dateutil import relativedelta, parser
from actions.parser import parse_time
from actions.dao import get_sum, get_count, get_tags

measurements = {
    "authorized_payments": "payment_authorized",
    "failed_payments": "payment_failed"
}


def get_time(tracker: Tracker,  dispatcher: CollectingDispatcher):
    time = next(tracker.get_latest_entity_values("time"), None)
    parsed_time = parse_time(
        tracker.latest_message) if time is not None else {}

    start_time = parsed_time.get("start_time", tracker.get_slot('start_time'))
    end_time = parsed_time.get("end_time", tracker.get_slot('end_time'))

    if start_time is None or end_time is None:
        dispatcher.utter_message(
            "Please specify a time or duration like today, last week, last friday to this monday etc")
        return (None, None, True)

    return (start_time, end_time, False)


def get_measurement(tracker:  Tracker,  dispatcher: CollectingDispatcher):
    measurement = next(tracker.get_latest_entity_values(
        "measurement"), tracker.get_slot('measurement'))

    if measurement not in ["authorized_payments", "failed_payments"]:
        dispatcher.utter_message(
            "Could not recognize metric: {}".format(measurement))
        return (None, True)

    return (measurement, False)


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

        tags = get_tags(measurements[measurement])

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

        method = next(
            tracker.get_latest_entity_values("method"), None)

        measurement, error = get_measurement(tracker, dispatcher)
        if error:
            return []

        start_time, end_time, error = get_time(tracker, dispatcher)
        if error:
            return []

        count = get_count(
            measurements[measurement], start_time, end_time, method)

        dispatcher.utter_message("Count: {}".format(count))
        return [SlotSet("start_time", start_time), SlotSet("end_time", end_time), SlotSet("measurement", measurement), SlotSet("mode", "count")]


class ActionSumMeasurements(Action):
    def name(self) -> Text:
        return "action_sum_measurements"

    def run(self,
            dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        method = next(
            tracker.get_latest_entity_values("method"), None)

        measurement, error = get_measurement(tracker, dispatcher)
        if error:
            return []

        start_time, end_time, error = get_time(tracker, dispatcher)
        if error:
            return []

        sum = get_sum(measurements[measurement],
                      start_time, end_time, method)

        dispatcher.utter_message("INR: {}".format(sum/100))
        return [SlotSet("start_time", start_time), SlotSet("end_time", end_time), SlotSet("measurement", measurement), SlotSet("mode", "sum")]


class ActionGetSuccessRate(Action):
    def name(self) -> Text:
        return "action_get_success_rate"

    def run(self,
            dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        method = next(
            tracker.get_latest_entity_values("method"), None)

        print(method)

        start_time, end_time, error = get_time(tracker, dispatcher)

        if error:
            return []

        authorized_count = get_count(
            "payment_authorized", start_time, end_time, method)

        failure_count = get_count(
            "payment_failed", start_time, end_time, method)

        success_rate = 0 if (authorized_count + failure_count ==
                             0) else authorized_count/(authorized_count + failure_count)

        dispatcher.utter_message(
            "Success Rate for : {:.2f}%".format(success_rate*100))
        return [SlotSet("start_time", start_time), SlotSet("end_time", end_time), SlotSet("mode", "success_rate")]
