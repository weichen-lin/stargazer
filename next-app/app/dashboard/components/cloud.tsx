import React, { useState } from 'react'
import { Text } from '@visx/text'
import { scaleLog } from '@visx/scale'
import Wordcloud from '@visx/wordcloud/lib/Wordcloud'

interface ExampleProps {
  width: number
  height: number
  showControls?: boolean
}

export interface WordData {
  text: string
  value: number
}

const colors = ['#143059', '#2F6B9A', '#82a6c2']

function getRotationDegree() {
  const rand = Math.random()
  const degree = rand > 0.5 ? 60 : -60
  return rand * degree
}

const words = [
  { text: '12321', value: 1234 },
  { text: 'qwe', value: 13 },
  { text: '123asdwq21', value: 3 },
  { text: 'asd', value: 12344 },
  { text: '12321', value: 123 },
  { text: '123sc21', value: 12222123 },
  { text: '12321', value: 11 },
  { text: '123svsdv21', value: 23 },
  { text: '123werwe21asdad', value: 13 },
  { text: '123wsvwe21', value: 1 },
  { text: '12werwe321', value: 3 },
]

const fontScale = scaleLog({
  domain: [Math.min(...words.map(w => w.value)), Math.max(...words.map(w => w.value))],
  range: [10, 100],
})
const fontSizeSetter = (datum: WordData) => fontScale(datum.value)

const fixedValueGenerator = () => 0.5

type SpiralType = 'archimedean' | 'rectangular'

export default function Example({ width, height, showControls }: ExampleProps) {
  const [spiralType, setSpiralType] = useState<SpiralType>('archimedean')
  const [withRotation, setWithRotation] = useState(false)

  return (
    <div className='flex flex-col select-none'>
      <Wordcloud
        words={words}
        width={width}
        height={height}
        fontSize={fontSizeSetter}
        font={'Impact'}
        padding={2}
        spiral={spiralType}
        rotate={withRotation ? getRotationDegree : 0}
        random={fixedValueGenerator}
      >
        {cloudWords =>
          cloudWords.map((w, i) => (
            <Text
              key={w.text}
              fill={colors[i % colors.length]}
              textAnchor={'middle'}
              transform={`translate(${w.x}, ${w.y}) rotate(${w.rotate})`}
              fontSize={w.size}
              fontFamily={w.font}
              onMouseEnter={e => {
                console.log(e)
              }}
            >
              {w.text}
            </Text>
          ))
        }
      </Wordcloud>
      {!showControls && (
        <div>
          <label>
            Spiral type &nbsp;
            <select onChange={e => setSpiralType(e.target.value as SpiralType)} value={spiralType}>
              <option key={'archimedean'} value={'archimedean'}>
                archimedean
              </option>
              <option key={'rectangular'} value={'rectangular'}>
                rectangular
              </option>
            </select>
          </label>
          <label>
            With rotation &nbsp;
            <input type='checkbox' checked={withRotation} onChange={() => setWithRotation(!withRotation)} />
          </label>
          <br />
        </div>
      )}
      {/* <style jsx>{`
        .wordcloud {
          display: flex;
          flex-direction: column;
          user-select: none;
        }
        .wordcloud svg {
          margin: 1rem 0;
          cursor: pointer;
        }

        .wordcloud label {
          display: inline-flex;
          align-items: center;
          font-size: 14px;
          margin-right: 8px;
        }
        .wordcloud textarea {
          min-height: 100px;
        }
      `}</style> */}
    </div>
  )
}
