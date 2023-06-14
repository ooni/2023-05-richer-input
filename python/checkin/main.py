#!/usr/bin/env python3

"""
Generates check-in v2 API responses.
"""

import argparse
import http.client
import json
import logging
import sys
from typing import Any, List


def log_json_message_to_stderr(prefix: str, message: Any):
    """Logs a JSON message to the stderr using the given prefix."""
    for line in json.dumps(message, indent=4).splitlines():
        print(prefix, line, file=sys.stderr)


def new_check_in_v2_request(probe_asn: str, probe_cc: str, only_categories: List[str]):
    """Generates a check-in v2 request."""
    return {
        "charging": True,
        "engine_name": "ooniprobe-engine",
        "engine_version": "0.1.0",
        "on_wifi": True,
        "platform": "linux",
        "probe_asn": probe_asn,
        "probe_cc": probe_cc,
        "run_type": "manual",
        "software_name": "miniooni",
        "software_version": "0.1.0",
        "web_connectivity": {
            "only_categories": only_categories,
        },
    }


def new_check_in_v1_request(v2_request):
    """Converts a check-in v2 API request to a check-in v1 API request."""
    return {
        "charging": v2_request["charging"],
        "on_wifi": v2_request["on_wifi"],
        "platform": v2_request["platform"],
        "probe_asn": v2_request["probe_asn"],
        "probe_cc": v2_request["probe_cc"],
        "run_type": v2_request["run_type"],
        "software_name": v2_request["software_name"],
        "software_version": v2_request["software_version"],
        "web_connectivity": {
            "category_codes": v2_request["web_connectivity"]["only_categories"],
        },
    }


def call_check_in_v1_api(request):
    """Calls the check-in v1 API."""
    raw_request = json.dumps(request).encode("utf-8")

    # create the HTTP connection
    http_conn = http.client.HTTPSConnection("api.ooni.io")

    # prepare the HTTP request
    headers = {
        "Content-Type": "application/json",
        "Content-Length": len(raw_request),
    }

    # send the HTTP request
    http_conn.request("POST", "/api/v1/check-in", body=raw_request, headers=headers)

    # Obtain the HTTP response
    http_response = http_conn.getresponse()

    # get the raw response body
    raw_response = http_response.read().decode("utf-8")

    # parse and return the check-in v1 response
    return json.loads(raw_response)


def geolocate():
    """Returns the current geolocation."""
    # That's curious: we need to omit the probe_asn and the probe_cc to obtain
    # the geolocation, which means this code will break if the probe tries to
    # include probe_asn and probe_cc into its request definition UNLESS we set
    # `json:"...,omitempty"`. Perhaps, we should make check-in v1 more robust
    # to ignore probe_asn and probe_cc when they're empty strings?
    request = {
        "charging": False,
        "on_wifi": False,
        "platform": "linux",
        # "probe_asn": "",
        # "probe_cc": "",
        "run_type": "timed",
        "software_name": "miniooni",
        "software_version": "0.1.0-dev",
        "web_connectivity": {"category_codes": ["MISC"]},
    }

    # We're abusing the /api/v1/check-in API to obtain our geolocation
    response = call_check_in_v1_api(request)

    return response["probe_asn"], response["probe_cc"]


def new_nettest_targets(nettest_name, nettest_entry):
    """Generates nettest targets for the given nettest."""
    if nettest_name == "web_connectivity":
        return nettest_entry["urls"]

    # TODO(bassosimone): we should evolve the richer input implementation
    # and therefore generate targets for other nettests.

    return []


def new_interpreter_script(v1_response):
    """Generates an interpreter script from the check-in v1 response."""
    # define the suites we generate by default
    suites = (
        ("websites", ("web_connectivity",)),
        ("im", ("facebook_messenger", "signal", "telegram", "whatsapp")),
        (
            "performance",
            (
                "dash",
                "ndt",
                "http_header_field_manipulation",
                "http_invalid_request_line",
            ),
        ),
        ("circumvention", ("psiphon", "tor")),
        (
            "experimental",
            (
                "dnscheck",
                "stun_reachability",
                "torsf",
                "vanilla_tor",
                "urlgetter",
            ),
        ),
    )

    # obtain the tests key of the response
    tests = v1_response["tests"]

    # generate the script for each suite
    script = []
    for suite_name, nettests in suites:
        # append the initial command setting the suite
        script.append(
            {
                "run_command": "ui/set_suite",
                "with_arguments": {"suite_name": suite_name},
            }
        )

        # compute how the progress bar should make progress
        if len(nettests) <= 0:
            continue
        increment = 1 / float(len(nettests))
        current_minimum = 0
        current_maximum = increment

        # generate an entry for each nettest
        for nettest_name in nettests:
            # obtain the nettest entry from the v1_response
            nettest_entry = tests.get(nettest_name)
            if nettest_entry:

                # position the progress bar correctly
                script.append(
                    {
                        "run_command": "ui/set_progress_bar_range",
                        "with_arguments": {
                            "initial_value": current_minimum,
                            "max_value": current_maximum,
                            "suite_name": suite_name,
                        },
                    }
                )

                # create the command running the nettest
                script.append(
                    {
                        "run_command": "nettest/run",
                        "with_arguments": {
                            "experimental_flags": {},
                            "nettest_name": nettest_name,
                            "report_id": nettest_entry["report_id"],
                            "suite_name": suite_name,
                            "targets": new_nettest_targets(nettest_name, nettest_entry),
                        },
                    }
                )

            else:
                logging.warning(f"cannot find {nettest_name} in v1 response")

            # move forward the progress bar
            current_minimum = current_maximum
            current_maximum += increment

        # make sure the progress bar is always at 100% after
        # we have run a given nettest suite
        script.append(
            {
                "run_command": "ui/set_progress_bar_value",
                "with_arguments": {
                    "suite_name": suite_name,
                    "value": 1,
                },
            }
        )

    return script


def new_check_in_v2_response(v1_response):
    """Converts a check-in v1 API response to a check-in v2 API response."""

    v2_response = {
        # add global configuration
        "config": {
            "test_helpers": v1_response["conf"]["test_helpers"],
        },
        # add the main measurement script
        "main_script": new_interpreter_script(v1_response),
        # add the fallback script that the probe should cache and only
        # use when it's not possible to contact the backend
        #
        # TODO(bassosimone): actually define and add the script
        #
        # TODO(bassosimone): add cache expiry time?
        "fallback_script": [],
        # add UTC time info
        "utc_time": v1_response["utc_time"],
        # add the version number
        "v": 2,
    }

    return v2_response


def main():
    # create the argument parser
    parser = argparse.ArgumentParser(
        prog="checkin",
        description="Generates check-in v2 API responses.",
    )

    # register the --only-category argument
    parser.add_argument(
        "--only-category",
        action="append",
        default=list(),
        help="only get URLs for a given category (can be specified multiple times)",
    )

    # parse the command line arguments
    args = parser.parse_args()

    # obtain the probe geolocation
    probe_asn, probe_cc = geolocate()

    # obtain the check-in v2 API request
    check_in_v2_request = new_check_in_v2_request(
        probe_asn, probe_cc, args.only_category
    )

    # log the check-in v2 API request
    log_json_message_to_stderr("v2request>", check_in_v2_request)

    # obtain the check-in v1 API request
    check_in_v1_request = new_check_in_v1_request(check_in_v2_request)

    # log the check-in v1 API request
    log_json_message_to_stderr("v1request>", check_in_v1_request)

    # call the check-in v1 API
    check_in_v1_response = call_check_in_v1_api(check_in_v1_request)

    # obtain the check-in v2 API response
    check_in_v2_response = new_check_in_v2_response(check_in_v1_response)

    # dump the v2 response to stdout
    json.dump(check_in_v2_response, sys.stdout)
    sys.stdout.write("\n")


if __name__ == "__main__":
    main()
