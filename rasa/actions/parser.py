from dateutil import relativedelta, parser
from typing import Dict, Text, Any, Optional


def parse_time(latest_message):
    timeinfo = next((x for x in latest_message.get(
        "entities") if x["entity"] == "time"), {}).get("additional_info", None)
    if timeinfo.get("type") == "interval":
        return close_interval_duckling_time(timeinfo)
    elif timeinfo.get("type") == "value":
        return make_interval_from_value_duckling_time(timeinfo)


def close_interval_duckling_time(
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
    timeinfo: Dict[Text, Any]
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
