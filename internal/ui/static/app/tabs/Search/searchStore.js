import { apiSearchTv } from "../../api.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { computed, defineStore, ref } from "../../vue.js"

export const useSearchStore = defineStore("search", () => {
  const loading = ref(false)

  const searchTerm = ref("")
  const page = ref(1)
  const pageTotal = ref(1)
  const resultsCnt = ref(0)
  const results = ref({})

  const onSearch = async (searchText, pg) => {
    if (loading.value) return

    if (!pg) pg = 1
    pg = Number(pg) || 1

    if (searchText === searchTerm.value && page.value === pg) {
      return
    }

    if (searchText === searchTerm.value && results.value[pg]) {
      page.value = pg
      return
    }

    try {
      loading.value = true
      page.value = pg

      const { data, err } = await apiSearchTv(searchText, pg)
      if (err) {
        throw err
      }

      if (searchTerm.value !== searchText) {
        searchTerm.value = searchText
        results.value = {}
      }
      pageTotal.value = data.TotalPages || 1
      resultsCnt.value = data.TotalResults || 0
      results.value = {
        ...results.value,
        [pg]: data.Results,
      }
    } catch (error) {
      console.error(error)
      notify(MsgType.Error, "Error", error.message)
    } finally {
      loading.value = false
    }
  }

  const getTotalPages = () => {
    return pageTotal.value
  }

  return {
    // data
    results,
    loading,
    getTotalPages,
    pageTotal: computed(() => pageTotal.value),
    page: computed(() => page.value),
    resultsCnt,

    // actions
    onSearch,
  }
})
