import {
    AnimatedAxis,
    AnimatedGrid,
    AnimatedLineSeries,
    GlyphSeries,
    Tooltip,
    XYChart,
} from "@visx/xychart";
import React, { useEffect, useRef, useState } from "react";
import { type ApiResult, get } from "./utils/api";
import { RenderTooltipParams } from "@visx/xychart/lib/components/Tooltip";

type Metrics = {
    id: string;
    route: string;
    method: string;
    statusCode: number;
    responseTime: number;
    createdOn: string;
    createdOnDate?: Date;
};

const unUsedcolors = [
    "#ff0000",
    "#00ff00",
    "#0000ff",
    "#ff00ff",
    "#00ffff",
    "#ffa500",
    "#800000",
    "#808000",
    "#008080",
    "#800080",
];

const usedColors = new Set<string>();

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

export function LiveGraph() {
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
    // TODO: handle the case when no line is shown(the graph disappears)

    return (
        <>
            <XYChart
                width={800}
                height={400}
                xScale={{ type: "time", domain: xAxisLimits }}
                yScale={{ type: "linear", domain: [0, 1000] }}
            >
                <AnimatedGrid columns={false} />
                <AnimatedAxis orientation="bottom" numTicks={4} />
                <AnimatedAxis orientation="left" />
                {liveChartData
                    .filter((l) => l.display)
                    .map((line) => {
                        return (
                            <GlyphSeries
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

                <Tooltip
                    snapTooltipToDatumX
                    snapTooltipToDatumY
                    renderTooltip={({
                        tooltipData,
                    }: RenderTooltipParams<Point>) => {
                        return (
                            <div>
                                <div>
                                    time:{" "}
                                    {tooltipData?.nearestDatum?.datum?.x &&
                                        new Date(
                                            tooltipData?.nearestDatum?.datum.x,
                                        ).toLocaleTimeString()}
                                </div>
                                <div>
                                    responseTime:{" "}
                                    {tooltipData?.nearestDatum?.datum?.y &&
                                        `${tooltipData?.nearestDatum?.datum?.y}ms`}
                                </div>
                            </div>
                        );
                    }}
                />
            </XYChart>
            <div className="flex flex-wrap gap-5">
                {liveChartData.map((ld) => {
                    return (
                        <div key={ld.route} className="flex gap-1 items-center">
                            <input
                                type="checkbox"
                                onChange={() =>
                                    handleRouteDisplayCheckboxClick(ld)
                                }
                                checked={ld.display}
                            />
                            <span style={{ color: ld.color }}>{ld.route}</span>
                        </div>
                    );
                })}
            </div>
        </>
    );
}

type RouteMetricsProps = {
    route: string;
    latestWindowInMinutes?: number;
};
export function RouteMetrics({
    route,
    latestWindowInMinutes = 10,
}: React.PropsWithChildren<RouteMetricsProps>) {
    const [metricsData, setMetricsData] = useState<Point[] | null>(null);

    // biome-ignore lint/correctness/useExhaustiveDependencies: this is componentDidMount.
    useEffect(() => {
        console.log(
            "route: ",
            route,
            "latestWindowInMinutes:",
            latestWindowInMinutes,
        );
        const fn = async () => {
            const queryParams = new URLSearchParams({
                route,
                latestWindowInMinutes: latestWindowInMinutes.toString(),
            }).toString();
            const res = await get<Metrics[], string>(
                `/api/metrics?${queryParams}`,
            );
            //console.log("metrics res: ", res);
            if (!res.ok || !res.ok.length) {
                return;
            }
            const data = res.ok.map((m) => {
                const p: Point = {
                    x: new Date(m.createdOn).getTime(),
                    y: m.responseTime,
                };
                return p;
            });
            setMetricsData(data);
        };
        fn();
    }, [route]);

    if (!metricsData) {
        return <h2>No metrics data for route: '{route}'</h2>;
    }

    return (
        <div className="flex flex-col items-center mt-4">
            <h2 className="font-bold">Metrics data for route: '{route}'</h2>
            <XYChart
                width={800}
                height={400}
                xScale={{ type: "time" }}
                yScale={{ type: "linear" }}
            >
                <AnimatedGrid columns={false} />
                <AnimatedAxis orientation="bottom" numTicks={4} />
                <AnimatedAxis orientation="left" />
                <GlyphSeries
                    key={`glyph-${route}`}
                    dataKey={`glyph-${route}`}
                    data={metricsData}
                    xAccessor={(d) => d?.x}
                    yAccessor={(d) => d.y}
                    renderGlyph={({ x, y }) => (
                        <circle
                            cx={x}
                            cy={y}
                            r={2} // Radius of the dots
                            fill="red" // Color of the dots
                        />
                    )}
                />
                <AnimatedLineSeries
                    key={`line-${route}`}
                    dataKey={`line-${route}`}
                    data={metricsData}
                    xAccessor={(d) => d?.x}
                    yAccessor={(d) => d.y}
                    stroke="red"
                />
                );
            </XYChart>
        </div>
    );
}
