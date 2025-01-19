import { useEffect, useRef, useState } from "react";
import { get, type ApiResult } from "./utils/api";
import {
    XYChart,
    AnimatedAxis,
    AnimatedGrid,
    AnimatedLineSeries,
    Tooltip,
} from "@visx/xychart";

type Metrics = {
    id: string;
    route: string;
    method: string;
    statusCode: number;
    responseTime: number;
    createdOn: string;
    createdOnDate?: Date;
};

//const generateData = (lineIndex: number) => {
//    return Array.from({ length: 10 }, (_, x) => ({
//        x,
//        y: Math.sin(x / 2 + lineIndex) * 10 + lineIndex * 2,
//    }));
//};

const unUsedcolors = [
    "#ff0000",
    "#00ff00",
    "#0000ff",
    "#ff00ff",
    "#00ffff",
    "#ffff00",
    "#800000",
    "#808000",
    "#008080",
    "#800080",
];

const usedColors = new Set<string>();

//const data = Array.from({ length: 10 }, (_, i) => ({
//    id: `Line ${i + 1}`,
//    data: generateData(i),
//    color: colors[i],
//}));

type Point = {
    x: Date;
    y: number;
};

type MetricsChartDataItem = {
    route: string;
    data: Point[];
    color: string;
};

const getNextUnusedColor = () => {
    const n = unUsedcolors.length;
    for (let i = 0; i < n; i++) {
        const color = unUsedcolors[i];
        if (!usedColors.has(color)) {
            usedColors.add(color);
            return color;
        }
    }
    return "red";
};

function App() {
    const [liveChartData, setLiveChartData] = useState<MetricsChartDataItem[]>(
        [],
    );
    const liveChartDataRef = useRef<MetricsChartDataItem[]>([]);
    const startWebsocket = () => {
        const ws = new WebSocket("/api/ws");
        ws.onmessage = (ev) => {
            try {
                const data = JSON.parse(ev.data) as ApiResult<Metrics, string>;
                if (data.ok) {
                    //console.log(data.ok);
                    const metricsData = data.ok;
                    const route = metricsData.route;
                    const chartData = {
                        x: new Date(metricsData.createdOn),
                        y: metricsData.responseTime,
                    };
                    const liveChartData = liveChartDataRef.current;

                    const existingRouteLiveChartData = liveChartData.find(
                        (d) => d.route === route,
                    );
                    if (!existingRouteLiveChartData) {
                        const newRouteLiveChartData = {
                            route,
                            data: [chartData],
                            color: getNextUnusedColor(),
                        };
                        // TODO: do some sorting so that the order of the lines are not changed.
                        const newLiveChartData = [
                            ...liveChartData,
                            newRouteLiveChartData,
                        ];
                        liveChartDataRef.current = newLiveChartData;
                    } else {
                        const newExistingRouteLiveChartData = {
                            ...existingRouteLiveChartData,
                        };
                        newExistingRouteLiveChartData.data = [
                            ...newExistingRouteLiveChartData.data,
                            chartData,
                        ];
                        liveChartDataRef.current = liveChartData.map((lc) => {
                            if (
                                lc.route !== newExistingRouteLiveChartData.route
                            ) {
                                return lc;
                            }
                            return newExistingRouteLiveChartData;
                        });
                    }
                }
            } catch (e) {
                console.log("error while receiving data: ", e);
            }
        };
    };

    // biome-ignore lint/correctness/useExhaustiveDependencies: want to treat this useEffect as componentDidMount
    useEffect(() => {
        const fn = async () => {
            const res = await get<Record<string, string>, string>(
                "/api/metrics",
            );
            console.log("here, ", res);
        };
        fn();
        startWebsocket();

        const timer = setInterval(() => {
            setLiveChartData(liveChartDataRef.current);
        }, 5000);

        return () => {
            clearInterval(timer);
        };
    }, []);

    return (
        <div className="flex flex-col p-1 items-center h-full w-full">
            <h1>Rocket tutor dashboard</h1>
            <div>
                <XYChart
                    width={800}
                    height={400}
                    xScale={{ type: "time" }}
                    yScale={{ type: "linear" }}
                >
                    <AnimatedGrid columns={false} />
                    <AnimatedAxis
                        orientation="bottom"
                        numTicks={5}
                        tickFormat={(value) =>
                            value.toLocaleTimeString("en-US", { hour12: false })
                        }
                    />
                    <AnimatedAxis orientation="left" />
                    {liveChartData.map((line) => (
                        <AnimatedLineSeries
                            key={line.route}
                            dataKey={line.route}
                            data={line.data}
                            xAccessor={(d) => d?.x}
                            yAccessor={(d) => d.y}
                            stroke={line.color}
                        />
                    ))}
                    {/*
                    <Tooltip
                        showVerticalCrosshair
                        snapTooltipToDatumX
                        snapTooltipToDatumY
                        renderTooltip={({ tooltipData }) => (
                            <div>
                                <div>
                                    <strong>
                                        {tooltipData?.nearestDatum?.key}
                                    </strong>
                                </div>
                                <div>
                                    x: {tooltipData?.nearestDatum?.datum?.x}
                                </div>
                                <div>
                                    y: {tooltipData?.nearestDatum?.datum?.y}
                                </div>
                            </div>
                        )}
                    />*/}
                </XYChart>
            </div>
            <div>
                {/*
                <XYChart
                    width={800}
                    height={400}
                    xScale={{ type: "linear" }}
                    yScale={{ type: "linear" }}
                >
                    <AnimatedGrid columns={false} numTicks={4} />
                    <AnimatedAxis orientation="bottom" />
                    <AnimatedAxis orientation="left" />
                    {data.map((line) => (
                        <AnimatedLineSeries
                            key={line.id}
                            dataKey={line.id}
                            data={line.data}
                            xAccessor={(d) => d.x}
                            yAccessor={(d) => d.y}
                            stroke={line.color}
                        />
                    ))}
                    <Tooltip
                        showVerticalCrosshair
                        snapTooltipToDatumX
                        snapTooltipToDatumY
                        renderTooltip={({ tooltipData }) => (
                            <div>
                                <div>
                                    <strong>
                                        {tooltipData?.nearestDatum?.key}
                                    </strong>
                                </div>
                                <div>
                                    x: {tooltipData?.nearestDatum?.datum?.x}
                                </div>
                                <div>
                                    y: {tooltipData?.nearestDatum?.datum?.y}
                                </div>
                            </div>
                        )}
                    />
                </XYChart>
                */}
            </div>
        </div>
    );
}

export default App;
