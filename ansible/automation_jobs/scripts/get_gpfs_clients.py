import sys
import requests


def get_gpfs_hosts(api_url):
    response = requests.get(api_url)

    if response.status_code == 200:
        data = response.json()
        return data['data']
    else:
        print(f"Error accessing API: {response.status_code}")
        return None


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python main.py <api_url>")
        sys.exit(1)
    url = 'https://devops.rsq.cn/api/public/gpfs/' + sys.argv[1]
    # url = 'http://127.0.0.1:8866/api/public/gpfs/' + sys.argv[1]
    result = get_gpfs_hosts(url)

    if result is not None:
        print(result)