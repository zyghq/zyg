import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { Area, AreaChart, Bar, BarChart, CartesianGrid, XAxis } from "recharts";

export const description = "A snapshot of the number of threads in Todo.";

const chartData = [
  { desktop: 186, month: "January" },
  { desktop: 305, month: "February" },
  { desktop: 237, month: "March" },
  { desktop: 73, month: "April" },
  { desktop: 209, month: "May" },
  { desktop: 214, month: "June" },
];

const chartConfig = {
  desktop: {
    color: "hsl(var(--chart-1))",
    label: "Desktop",
  },
} satisfies ChartConfig;

export function QueueSize({ className }: { className?: string }) {
  return (
    <ChartContainer className={className} config={chartConfig}>
      <AreaChart
        accessibilityLayer
        data={chartData}
        margin={{
          left: 12,
          right: 12,
        }}
      >
        <CartesianGrid vertical={false} />
        <XAxis
          axisLine={false}
          dataKey="month"
          tickFormatter={(value) => value.slice(0, 3)}
          tickLine={false}
          tickMargin={8}
        />
        <ChartTooltip
          content={<ChartTooltipContent hideLabel indicator="dot" />}
          cursor={false}
        />
        <Area
          dataKey="desktop"
          fill="var(--color-desktop)"
          fillOpacity={0.4}
          stroke="var(--color-desktop)"
          type="linear"
        />
      </AreaChart>
    </ChartContainer>
  );
}

const volumeData = [
  { desktop: 186, month: "January" },
  { desktop: 305, month: "February" },
  { desktop: 237, month: "March" },
  { desktop: 73, month: "April" },
  { desktop: 209, month: "May" },
  { desktop: 214, month: "June" },
];

export function Volume({ className }: { className?: string }) {
  return (
    <ChartContainer className={className} config={chartConfig}>
      <BarChart accessibilityLayer data={volumeData}>
        <CartesianGrid vertical={false} />
        <XAxis
          axisLine={false}
          dataKey="month"
          tickFormatter={(value) => value.slice(0, 3)}
          tickLine={false}
          tickMargin={10}
        />
        <ChartTooltip
          content={<ChartTooltipContent hideLabel />}
          cursor={false}
        />
        <Bar dataKey="desktop" fill="var(--color-desktop)" radius={8} />
      </BarChart>
    </ChartContainer>
  );
}
