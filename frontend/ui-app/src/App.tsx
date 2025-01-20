import { useEffect, useRef, useState } from "react";
import { get, type ApiResult } from "./utils/api";
import {
    XYChart,
    AnimatedAxis,
    AnimatedGrid,
    AnimatedLineSeries,
    GlyphSeries,
    Tooltip,
} from "@visx/xychart";

import { curveBasis, curveMonotoneX, curveNatural } from "@visx/curve";
import { Fragment } from "react";

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
    x: number;
    y: number;
};

type MetricsChartDataItem = {
    route: string;
    data: Point[];
    color: string;
    display: boolean;
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

const getXAxisLimits = () => {
    const time = new Date().getTime();
    return [time - 10 * 60 * 1000, time];
};

function App() {
    const [liveChartData, setLiveChartData] = useState<MetricsChartDataItem[]>(
        [],
    );
    const [xAxisLimits, setXAxisLimits] = useState(getXAxisLimits);
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
                        x: new Date(metricsData.createdOn).getTime(),
                        y: metricsData.responseTime,
                    };
                    const liveChartData = liveChartDataRef.current;

                    const existingRouteLiveChartData = liveChartData.find(
                        (d) => d.route === route,
                    );
                    if (!existingRouteLiveChartData) {
                        const newRouteLiveChartData: MetricsChartDataItem = {
                            route,
                            data: [chartData],
                            color: getNextUnusedColor(),
                            display: true,
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

    const chartUpdate = () => {
        const xLimits = getXAxisLimits();
        setXAxisLimits(xLimits);

        setLiveChartData(
            liveChartDataRef.current.map((ld) => {
                return {
                    ...ld,
                    data: ld.data.filter((d) => d.x >= xLimits[0]),
                };
            }),
        );
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

        const timer = setInterval(chartUpdate, 1000);

        return () => {
            clearInterval(timer);
        };
    }, []);

    const handleRouteDisplayCheckboxClick = (ld: MetricsChartDataItem) => {
        setLiveChartData(
            liveChartDataRef.current.map((d) => {
                if (d.route !== ld.route) {
                    return d;
                }
                d.display = !d.display;
                const newLd: MetricsChartDataItem = {
                    ...d,
                };
                return newLd;
            }),
        );
    };

    return (
        <div className="flex flex-col p-1 items-center h-full w-full">
            <h1>Rocket tutor dashboard</h1>
            <div>
                <XYChart
                    width={800}
                    height={400}
                    xScale={{ type: "time", domain: xAxisLimits }}
                    yScale={{ type: "linear", domain: [0, 1000] }}
                >
                    <AnimatedGrid columns={false} />
                    <AnimatedAxis
                        orientation="bottom"
                        numTicks={4}
                        //tickFormat={(value) =>
                        //    value.toLocaleTimeString("en-US", { hour12: false })
                        //}
                    />
                    <AnimatedAxis orientation="left" />
                    {liveChartData
                        .filter((l) => l.display)
                        .map((line) => {
                            return (
                                <AnimatedLineSeries
                                    key={`line-${line.route}`}
                                    dataKey={`line-${line.route}`}
                                    data={line.data}
                                    xAccessor={(d) => d?.x}
                                    yAccessor={(d) => d.y}
                                    stroke={line.color}
                                />
                            );
                        })}
                    {liveChartData
                        .filter((l) => l.display)
                        .map((line) => {
                            return (
                                <GlyphSeries
                                    //curve={curveMonotoneX}
                                    key={`glyph-${line.route}`}
                                    dataKey={`glyph-${line.route}`}
                                    data={line.data}
                                    xAccessor={(d) => d?.x}
                                    yAccessor={(d) => d.y}
                                    renderGlyph={({ x, y }) => (
                                        <circle
                                            cx={x}
                                            cy={y}
                                            r={2} // Radius of the dots
                                            fill={line.color} // Color of the dots
                                        />
                                    )}
                                />
                            );
                        })}
                    <Tooltip
                        snapTooltipToDatumX
                        snapTooltipToDatumY
                        renderTooltip={({ tooltipData }) => {
                            return (
                                <div>
                                    <div>
                                        <strong>
                                            {tooltipData?.nearestDatum?.key}
                                        </strong>
                                    </div>
                                    <div>
                                        time:{" "}
                                        {/* {tooltipData?.nearestDatum?.datum?.x?.toLocaleTimeString()}*/}
                                    </div>
                                    <div>
                                        responseTime:{" "}
                                        {tooltipData?.nearestDatum?.datum?.y}
                                    </div>
                                </div>
                            );
                        }}
                    />
                </XYChart>
            </div>
            <div className="flex flex-col">
                {liveChartData.map((ld) => {
                    return (
                        <div key={ld.route} className="flex">
                            <input
                                type="checkbox"
                                onChange={() =>
                                    handleRouteDisplayCheckboxClick(ld)
                                }
                                checked={ld.display}
                            />
                            <span>{ld.route}</span>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

export default App;
