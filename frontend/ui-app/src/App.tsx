import { LiveGraph, RouteMetrics } from "./Graphs";

function App() {
    // TODO: handle the case when no line is shown(the graph disappears)

    return (
        <div className="flex flex-col p-1 items-center h-full w-[80%]">
            <h1>Rocket tutor dashboard</h1>
            <LiveGraph />
            <RouteMetrics route="/metrics" />
            <div className="w-full"></div>
        </div>
    );
}

export default App;
