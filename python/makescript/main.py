#!/usr/bin/env python3

"""
Given the check-in v2 response, this command generates the
corresponding interpreter script from the "main" script.
"""

import json
import sys

def main():
    with open(sys.argv[1], "r") as filep:
        v2_response = json.load(filep)
        script = {
            "config": v2_response["config"],
            "commands": v2_response["main_script"],
            "v": v2_response["v"],
        }
        json.dump(script, sys.stdout)
        sys.stdout.write("\n")

if __name__ == "__main__":
    main()
