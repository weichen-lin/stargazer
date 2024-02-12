import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'
import { useSearchParams } from 'next/navigation'

export default function FixPagination(props: { total: number }) {
  const { total } = props
  const searchParams = useSearchParams()
  const current = parseInt(searchParams.get('p') ?? '1')

  return <Pagination>{renderPagination(total, current)}</Pagination>
}

const renderPagination = (total: number, current: number) => {
  const totalPage = Math.ceil(total / 20)

  return (
    <PaginationContent>
      {current > 1 && (
        <PaginationItem>
          <PaginationPrevious href={`/stars?p=${current - 1}`} />
        </PaginationItem>
      )}
      {current - 2 > 0 && (
        <PaginationItem>
          <PaginationLink href={`/stars?p=1`}>1</PaginationLink>
        </PaginationItem>
      )}
      {current - 2 > 1 && (
        <PaginationItem>
          <PaginationEllipsis />
        </PaginationItem>
      )}
      {current - 1 > 0 && (
        <PaginationItem>
          <PaginationLink href={`/stars?p=${current - 1}`}>{current - 1}</PaginationLink>
        </PaginationItem>
      )}
      <PaginationItem>
        <PaginationLink isActive>{current}</PaginationLink>
      </PaginationItem>
      {current + 1 < totalPage && (
        <PaginationItem>
          <PaginationLink href={`/stars?p=${current + 1}`}>{current + 1}</PaginationLink>
        </PaginationItem>
      )}
      {current + 2 < totalPage && (
        <PaginationItem>
          <PaginationEllipsis />
        </PaginationItem>
      )}
      {current + 1 == totalPage && (
        <PaginationItem>
          <PaginationLink href={`/stars?p=${totalPage}`}>{totalPage}</PaginationLink>
        </PaginationItem>
      )}
      {current + 2 <= totalPage && (
        <PaginationItem>
          <PaginationLink href={`/stars?p=${totalPage}`}>{totalPage}</PaginationLink>
        </PaginationItem>
      )}
      {current < totalPage && (
        <PaginationItem>
          <PaginationNext href={`/stars?p=${current + 1}`} />
        </PaginationItem>
      )}
    </PaginationContent>
  )
}
