import aiohttp
import asyncio
import json


async def create_monitor():
    project_id = "3255"
    token = "OEkftWB6p3JMXu3MVw9LhA"
    url = f"https://api2.uptrace.dev/internal/v1/projects/{project_id}/monitors"

    with open("monitor.json", "r") as f:
        data = json.load(f)

    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}

    print("\n===== REQUEST INFO =====")
    print(f"URL: {url}")
    print("Headers:")
    for k, v in headers.items():
        print(f"  {k}: {v}")

    print("\n===== REQUEST BODY =====")
    print(json.dumps(data, indent=2))

    async with aiohttp.ClientSession() as session:
        try:
            async with session.post(url, json=data, headers=headers) as response:
                print("===== RESPONSE INFO =====")
                print(f"Status Code: {response.status}")
                print(f"OK?: {response.ok}")
                print(f"Reason: {response.reason}")
                print(f"URL: {response.url}")
                print(f"Method: {response.method}")
                print(f"Content-Type: {response.content_type}")
                print(f"Charset: {response.charset}")
                print(f"Content-Encoding: {response.get_encoding()}")
                print("Headers:")
                for k, v in response.headers.items():
                    print(f"  {k}: {v}")

                print("\n===== RAW RESPONSE BODY =====")
                raw_body = await response.read()
                print(raw_body.decode("utf-8", errors="replace"))

                print("\n===== JSON (if possible) =====")
                try:
                    json_body = await response.json()
                    print(json.dumps(json_body, indent=2))
                except Exception as e:
                    print(f"Could not parse JSON: {e}")

        except aiohttp.ClientError as e:
            print("\n===== CLIENT ERROR =====")
            print(f"{type(e).__name__}: {e}")


# Run the async function
asyncio.run(create_monitor())
