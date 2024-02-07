import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'

export default function FixPagination(props: { total: number; current: number; page: string }) {
  const { total, current, page } = props
  return <Pagination>{renderPagination(total, current, page)}</Pagination>
}

const renderPagination = (total: number, current: number, page: string) => {
  const totalPage = Math.ceil(total / 20)

  return (
    <PaginationContent>
      {current > 1 && (
        <PaginationItem>
          <PaginationPrevious href={`${page}${current - 1}`} />
        </PaginationItem>
      )}
      {current - 2 > 0 && (
        <PaginationItem>
          <PaginationLink href={`${page}1`}>1</PaginationLink>
        </PaginationItem>
      )}
      {current - 2 > 1 && (
        <PaginationItem>
          <PaginationEllipsis />
        </PaginationItem>
      )}
      {current - 1 > 0 && (
        <PaginationItem>
          <PaginationLink href={`${page}${current - 1}`}>{current - 1}</PaginationLink>
        </PaginationItem>
      )}
      <PaginationItem>
        <PaginationLink isActive>{current}</PaginationLink>
      </PaginationItem>
      {current + 1 < totalPage && (
        <PaginationItem>
          <PaginationLink href={`${page}${current + 1}`}>{current + 1}</PaginationLink>
        </PaginationItem>
      )}
      {current + 2 < totalPage && (
        <PaginationItem>
          <PaginationEllipsis />
        </PaginationItem>
      )}
      {current + 1 == totalPage && (
        <PaginationItem>
          <PaginationLink href={`${page}${totalPage}`}>{totalPage}</PaginationLink>
        </PaginationItem>
      )}
      {current + 2 <= totalPage && (
        <PaginationItem>
          <PaginationLink href={`${page}${totalPage}`}>{totalPage}</PaginationLink>
        </PaginationItem>
      )}
      {current < totalPage && (
        <PaginationItem>
          <PaginationNext href={`${page}${current + 1}`} />
        </PaginationItem>
      )}
    </PaginationContent>
  )
}
