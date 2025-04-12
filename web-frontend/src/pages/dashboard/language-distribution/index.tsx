import { useMemo } from 'react';
import { Label, Pie, PieChart } from 'recharts';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import { colorConfig, getLanguageColor } from './config';
import type { Language } from './config';
import { PieChart as PieChartIcon, Plus } from 'lucide-react';
import { MagicCard } from '@/components/shared/magic-card';

const fakeData = [
  {
    language: 'TypeScript',
    count: 264,
  },
  {
    language: 'Go',
    count: 150,
  },
  {
    language: 'Python',
    count: 121,
  },
  {
    language: 'JavaScript',
    count: 105,
  },
  {
    language: 'Unknown',
    count: 36,
  },
  {
    language: 'Rust',
    count: 28,
  },
  {
    language: 'HTML',
    count: 12,
  },
  {
    language: 'C',
    count: 12,
  },
  {
    language: 'Java',
    count: 12,
  },
  {
    language: 'C++',
    count: 11,
  },
  {
    language: 'Shell',
    count: 11,
  },
  {
    language: 'PHP',
    count: 4,
  },
  {
    language: 'Swift',
    count: 4,
  },
  {
    language: 'Clojure',
    count: 3,
  },
  {
    language: 'MDX',
    count: 3,
  },
  {
    language: 'Jupyter Notebook',
    count: 3,
  },
  {
    language: 'Markdown',
    count: 3,
  },
  {
    language: 'Vue',
    count: 3,
  },
  {
    language: 'Dart',
    count: 2,
  },
  {
    language: 'Dockerfile',
    count: 2,
  },
  {
    language: 'OCaml',
    count: 1,
  },
  {
    language: 'SCSS',
    count: 1,
  },
  {
    language: 'HCL',
    count: 1,
  },
  {
    language: 'Ruby',
    count: 1,
  },
  {
    language: 'SVG',
    count: 1,
  },
  {
    language: 'Lua',
    count: 1,
  },
  {
    language: 'Jinja',
    count: 1,
  },
  {
    language: 'CSS',
    count: 1,
  },
  {
    language: 'Pug',
    count: 1,
  },
  {
    language: 'Zig',
    count: 1,
  },
  {
    language: 'Astro',
    count: 1,
  },
];

export default function LanguageDistribution() {
  const totalStars = useMemo(() => {
    return fakeData ? fakeData.reduce((acc, curr) => acc + curr.count, 0) : 0;
  }, [fakeData]);

  const withColor = fakeData
    ? fakeData.map((e) => ({
        language: e.language,
        count: e.count,
        fill: getLanguageColor(e.language as Language),
      }))
    : [];

  return (
    <Card className='flex flex-col h-[320px] w-full max-w-[380px] md:max-w-none'>
      <CardHeader className='items-center pb-0 gap-y-1'>
        <CardTitle className='text-xl'>Language Distribution</CardTitle>
      </CardHeader>
      <CardContent className='flex-1'>
        {fakeData.length > 0 && (
          <ChartContainer
            className='mx-auto aspect-square max-h-[250px] py-4'
            config={colorConfig}
          >
            <PieChart>
              <ChartTooltip
                cursor={false}
                content={<ChartTooltipContent hideLabel />}
              />
              <Pie
                data={withColor}
                dataKey='count'
                nameKey='language'
                innerRadius={60}
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
      </CardContent>
    </Card>
  );
}

const Loading = () => {
  return (
    <div className='flex flex-col items-center justify-center py-1 gap-y-4 pb-16'>
      <div className='w-[150px] h-[150px] rounded-full bg-slate-200 animate-pulse'></div>
    </div>
  );
};

const EmptyContent = () => {
  return (
    <div className='w-full flex flex-col items-center justify-center my-4'>
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
