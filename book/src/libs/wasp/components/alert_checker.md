# WASP - AlertChecker

Basically, if you have a dashboard you can either assert if any alert fired during a time range or if any group fired at all. Groups are defined by using label key `requirement_name.

For built-in alerts, we create a new row on the dashboard for each alert.
For the custom ones, we don't do it, so that we don't overload the dashboard.