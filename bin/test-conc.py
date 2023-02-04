import requests
import threading
import time

BASE_URL = "http://localhost:8080"
FILENAME = "9ffa0e25-7391-485d-a75a-58627d4b46a3.sh"


def sendrequest():
    r = requests.get(BASE_URL+"/file/?filename="+FILENAME)

    if r.status_code != 200:
        print("Response: HTTP %d" % r.status_code)

    time.sleep(1)


def main():
    maximum_requests = 32
    for _ in range(maximum_requests):
        th = threading.Thread(target=sendrequest)
        th.start()


main()
