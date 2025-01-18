import { useEffect } from "react";
import { get } from "./utils/api";

function App() {
    const startWebsocket = () => {
        const ws = new WebSocket("/api/ws");
        ws.onmessage = (ev) => {
            console.log(ev);
        };

        setTimeout(() => {
            ws.close();
        }, 4000);
    };

    // biome-ignore lint/correctness/useExhaustiveDependencies: want to treat this useEffect as componentDidMount
    useEffect(() => {
        const fn = async () => {
            console.log("hello, world!!");
            const res = await get<Record<string, string>, string>(
                "/api/metrics",
            );
            console.log("here, ", res);
        };
        fn();
        startWebsocket();
    }, []);
    return (
        <div>
            {" "}
            <h1>Rocket tutor dashboard</h1>
        </div>
    );
}

export default App;
