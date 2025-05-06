import { useMemo, useCallback } from 'react';
import { Label, Pie, PieChart } from 'recharts';
import { Card, CardContent } from '@/components/ui/card';
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import { colorConfig, getLanguageColor } from './config';
import type { Language } from './config';
import { PieChart as PieChartIcon, Plus } from 'lucide-react';

import useLanguageDistribution from '@/apis/repository/useLanguageDistribution';

export default function LanguageDistribution() {
  const { data, isLoading } = useLanguageDistribution();

  const totalStars = useMemo(() => {
    return data ? data.reduce((acc, curr) => acc + curr.count, 0) : 0;
  }, [data]);

  const withColor = data
    ? data.map((e) => ({
        language: e.language === '' ? 'Unknown' : e.language,
        count: e.count,
        fill: getLanguageColor(e.language as Language),
      }))
    : [];

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none border-none'>
      <CardContent className='flex items-center justify-center w-full h-full'>
        {isLoading && <Loading />}
        {!isLoading && data && data.length > 0 && (
          <ChartContainer config={colorConfig} className='w-full h-full'>
            <PieChart>
              <ChartTooltip
                cursor={false}
                content={<ChartTooltipContent hideLabel />}
              />
              <Pie
                data={withColor}
                dataKey='count'
                nameKey='language'
                innerRadius={70}
                strokeWidth={5}
              >
                <Label
                  content={({ viewBox }) => {
                    if (viewBox && 'cx' in viewBox && 'cy' in viewBox) {
                      return (
                        <text
                          x={viewBox.cx}
                          y={viewBox.cy}
                          textAnchor='middle'
                          dominantBaseline='middle'
                        >
                          <tspan
                            x={viewBox.cx}
                            y={viewBox.cy}
                            className='fill-foreground text-3xl font-bold'
                          >
                            {totalStars.toLocaleString()}
                          </tspan>
                          <tspan
                            x={viewBox.cx}
                            y={(viewBox.cy || 0) + 24}
                            className='fill-muted-foreground'
                          >
                            Stars
                          </tspan>
                        </text>
                      );
                    }
                  }}
                />
              </Pie>
            </PieChart>
          </ChartContainer>
        )}
        {!isLoading && !data && <EmptyContent />}
      </CardContent>
    </Card>
  );
}

const Loading = () => {
  return (
    <div className='w-[200px] h-[200px] rounded-full bg-slate-200 animate-pulse'></div>
  );
};

const EmptyContent = () => {
  return (
    <div className='w-full flex flex-col items-center justify-center'>
      <div className='w-32 h-32 relative'>
        <PieChartIcon className='w-full h-full text-gray-200' />
        <div className='absolute inset-0 flex items-center justify-center'>
          <Plus className='w-8 h-8 text-gray-400' />
        </div>
      </div>
      <p className='text-center text-gray-500 mb-4'>No data yet.</p>
    </div>
  );
};
