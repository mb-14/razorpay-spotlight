
from influxdb import InfluxDBClient

client = InfluxDBClient('localhost', 8086, '', '', 'rzpftx')


def get_sum(measurement, start_time, end_time, method=None):
    filter = "and method = '{}'".format(
        method) if (method is not None) else ""
    query = "SELECT SUM(amount) as amount FROM {} where time >= '{}' and time <= '{}' {}".format(
            measurement, start_time, end_time, filter)

    print(query)

    result = client.query(query)
    points = list(result.get_points())
    return points[0].get("amount", 0) if (len(points) == 1) else 0


def get_count(measurement, start_time, end_time, method=None):
    filter = "and method = '{}'".format(
        method) if (method is not None) else ""
    query = "SELECT COUNT(amount) as count_amount FROM {} where time >= '{}' and time <= '{}' {}".format(
            measurement, start_time, end_time, filter)

    print(query)

    result = client.query(query)
    points = list(result.get_points())
    return points[0].get("count_amount", 0) if (len(points) == 1) else 0


def get_tags(measurement):
    result = client.query(
        "SHOW TAG KEYS FROM {}".format(measurement))
    points = list(result.get_points())
    return map(lambda x: x["tagKey"], points)
