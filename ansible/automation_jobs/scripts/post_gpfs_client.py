import sys
import requests


def send_post_request(api_url, data_value):
    data = {"data": data_value}
    response = requests.post(api_url, json=data)

    if response.status_code == 200:
        print("POST request successful")
        print("Response:", response.json())
    else:
        print(f"Error: {response.status_code}")
        print("Response:", response.text)


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <api_url> <data_value>")
        sys.exit(1)

    api_url = 'https://devops.rsq.cn/api/public/gpfs/' + sys.argv[1]
    # api_url = 'http://127.0.0.1:8866/api/public/gpfs/' + sys.argv[1]
    data_value = sys.argv[2]

    send_post_request(api_url, data_value)