import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination'

export default function CollectionPagination(props: { total: number }) {
  const { total } = props

  return <Pagination>{renderPagination(total, 1, () => {})}</Pagination>
}

const renderPagination = (total: number, current: number, search: (page: number) => void) => {
  const totalPage = Math.ceil(total / 20)

  return (
    <PaginationContent>
      {current > 1 && (
        <PaginationItem>
          <PaginationPrevious onClick={() => search(current - 1)} />
        </PaginationItem>
      )}
      {current - 2 > 0 && (
        <PaginationItem>
          <PaginationLink onClick={() => search(1)}>1</PaginationLink>
        </PaginationItem>
      )}
      {current - 2 > 1 && (
        <PaginationItem>
          <PaginationEllipsis />
        </PaginationItem>
      )}
      {current - 1 > 0 && (
        <PaginationItem>
          <PaginationLink onClick={() => search(current - 1)}>{current - 1}</PaginationLink>
        </PaginationItem>
      )}
      <PaginationItem>
        <PaginationLink isActive>{current}</PaginationLink>
      </PaginationItem>
      {current + 1 < totalPage && (
        <PaginationItem>
          <PaginationLink onClick={() => search(current + 1)}>{current + 1}</PaginationLink>
        </PaginationItem>
      )}
      {current + 2 < totalPage && (
        <PaginationItem>
          <PaginationEllipsis />
        </PaginationItem>
      )}
      {current + 1 == totalPage && (
        <PaginationItem>
          <PaginationLink onClick={() => search(totalPage)}>{totalPage}</PaginationLink>
        </PaginationItem>
      )}
      {current + 2 <= totalPage && (
        <PaginationItem>
          <PaginationLink onClick={() => search(totalPage)}>{totalPage}</PaginationLink>
        </PaginationItem>
      )}
      {current < totalPage && (
        <PaginationItem>
          <PaginationNext onClick={() => search(current + 1)} />
        </PaginationItem>
      )}
    </PaginationContent>
  )
}
