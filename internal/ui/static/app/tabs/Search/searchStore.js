import { apiSearchTv } from "../../api.js"
import { notify } from "../../components/Notify/Notify.js"
import { defineStore, ref, useRouter } from "../../vue.js"

export const useSearchStore = defineStore("search", () => {
  const loading = ref(false)

  let searchTerm = ""
  let page = 1

  var pageTotal = 1
  const resultsCnt = ref(0)
  const results = ref({})

  const onSearch = async (searchText, pg) => {
    if (loading.value) return

    if (searchText === searchTerm && page === pg) {
      return
    }

    if (searchText === searchTerm && results.value[pg]) {
      page = pg
      return
    }

    try {
      loading.value = true
      if (!pg) {
        pg = 1
      }
      page = pg

      const { data, err } = await apiSearchTv(searchText, pg)
      if (err) {
        throw err
      }

      if (searchTerm !== searchText) {
        searchTerm = searchText
        results.value = {}
      }
      pageTotal = data.TotalPages
      resultsCnt.value = data.TotalResults
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
    return pageTotal
  }

  return {
    // data
    results,
    loading,
    getTotalPages,
    pageTotal,
    resultsCnt,

    // actions
    onSearch,
  }
})
