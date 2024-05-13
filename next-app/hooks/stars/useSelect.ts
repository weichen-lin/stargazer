import { create } from 'zustand'
import { Option } from '@/components/ui/multiple-selector'

interface SelectState {
  selected: Option[]
  setSelected: (selected: Option[]) => void
}

const useSelect = create<SelectState>(set => ({
  selected: [],
  setSelected: selected => set({ selected }),
}))

export default useSelect
