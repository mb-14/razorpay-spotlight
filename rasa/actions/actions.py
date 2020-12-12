# This files contains your custom actions which can be used to run
# custom Python code.
#
# See this guide on how to implement these action:
# https://rasa.com/docs/rasa/custom-actions

from rasa_sdk import Action, Tracker
from rasa_sdk.executor import CollectingDispatcher
from typing import Dict, Text, Any, List, Union, Optional
from influxdb import InfluxDBClient

client = InfluxDBClient('localhost', 8086, '', '', 'rzpftx')

measurements = {
       "payments": "payment_authorized",
       "failed_payments": "payment_failed"
       }

class ActionListMetrics(Action):
   def name(self) -> Text:
      return "action_list_metrics"

   def run(self,
           dispatcher: CollectingDispatcher,
           tracker: Tracker,
           domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

      result = ["payment_authorized (count/volume)", "payment_failure (count/volume)", "success_rate"]
    
      dispatcher.utter_message("The following metrics are available to query from: {}".format(", ".join(result)))
      return []


class ActionListDimensions(Action):

   def name(self) -> Text:
      return "action_list_dimensions"

   def run(self,
           dispatcher: CollectingDispatcher,
           tracker: Tracker,
           domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

      measurement = next(tracker.get_latest_entity_values("measurement"), None)
      if measurement not in ["payments", "failed_payments"]:
        dispatcher.utter_message("Could not recognize metric: {}".format(measurement))
        return []

      result = client.query("SHOW TAG KEYS FROM {}".format(measurements[measurement]))
      points = list(result.get_points())
      tags = map(lambda x: x["tagKey"], points)
      dispatcher.utter_message("The following dimensions are available for {}: {}".format(measurement, ", ".join(tags)))
      return []
