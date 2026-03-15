import json
import sys
import os
import asyncio
import websockets

token = os.getenv("HASS_TOKEN")
server = os.getenv("HASS_SERVER")

if not token or not server:
    print("Fehler: HASS_TOKEN und HASS_SERVER müssen gesetzt sein.")
    sys.exit(1)

# Umwandlung von http zu ws
ws_url = server.replace("http", "ws") + "/api/websocket"

async def listen():
    try:
        async with websockets.connect(ws_url) as websocket:
            # 1. Authentifizierung
            await websocket.send(json.dumps({
                "type": "auth",
                "access_token": token
            }))
            
            auth_resp = json.loads(await websocket.recv())
            if auth_resp["type"] != "auth_ok":
                print(f"Auth Fehler: {auth_resp}")
                return

            # 2. Auf Events abonnieren
            await websocket.send(json.dumps({
                "id": 1,
                "type": "subscribe_events",
                "event_type": "mobile_app_notification_action"
            }))
            
            sub_resp = json.loads(await websocket.recv())
            print(f"Lausche auf Action-Events...")

            # 3. Auf Nachricht warten
            while True:
                msg = json.loads(await websocket.recv())
                if msg["type"] == "event":
                    action = msg["event"]["data"]["action"]
                    print(f"ANTWORT ERHALTEN: {action}")
                    break
    except Exception as e:
        print(f"Fehler: {e}")

if __name__ == "__main__":
    asyncio.run(listen())
