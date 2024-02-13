import streamlit as st
import aiohttp
import asyncio
import base64
import streamlit.components.v1 as components

WS_CONN = "ws://localhost:9090/generate"

async def startDownload():
    async with aiohttp.ClientSession(trust_env = True) as session:
        async with session.ws_connect(WS_CONN) as websocket:
            await websocket.send_json(data={
                "filename": "test.csv"
            })
            message = await websocket.receive()
            b64 = base64.b64encode(message.data.encode()).decode()
            dl_link = f"""
                <a href="data:text/csv;base64,{b64}" id="downloadFile" download="test.csv">
                <script>
                    document.getElementById('downloadFile').click();
                </script>
            """
            components.html(
                dl_link,
                height=0,
            )

def download():
    asyncio.run(startDownload())

st.title("Metrics UI")
st.button("Process", on_click=download)