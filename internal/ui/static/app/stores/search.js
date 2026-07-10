import { ENDPOINT } from "../utils.js";
import { defineStore, ref } from "../vue.js";

const searchTv = async (search, page) => {
    let url = `${ENDPOINT.SEARCH_SERIES()}?q=${search}&p=${page}`
    try {
        const response = await fetch(url);
        const result = await response.json();
        return {
            data: result,
            err: null
        }
    } catch (error) {
        console.error('Error fetching music data:', error);
        return {
            data: result,
            err: error
        }
    }
}

export const useSearchStore = defineStore('search', () => {
    const loading = ref(false)

    let searchTerm = ''

    const pageCur = ref(1)
    const pageTotal = ref(1)
    const resultsCnt = ref(0)
    const results = ref({ 1: [] })

    const onSearch = async (searchText) => {
        
        if (loading.value || searchTerm === searchText) return
    
        try {
            loading.value = true
            const { data, err } = await searchTv(searchText, 1)
            if (err) {
                throw err
            }

            searchTerm = searchText
            pageCur.value = 1
            pageTotal.value = data.TotalPages
            resultsCnt.value = data.TotalResults
            results.value = {
                1: data.Results
            }

        } catch (error) {
            console.error(error)
        } finally {
            loading.value = false
        }
    
    }

    const onPrvBtn = async () => {
   
        if (pageCur.value === 1) {
            return
        }
        
        if (results.value[pageCur.value - 1]) {
            pageCur.value -= 1
            return
        }
        
        
        try {
            loading.value = true
            const { data, err } = await searchTv(searchTerm, pageCur.value - 1)
            if (err) {
                throw err
            }
    
    
            pageTotal.value = data.TotalPages
            resultsCnt.value = data.TotalResults
            results.value = {
                ...results.value,
                [pageCur.value - 1]: data.Results
            }
            pageCur.value -= 1
    
        } catch (error) {
            console.error(error)
        } finally {
            loading.value = false
        }
    
    }
    
    const onNxtBtn = async () => {
   
        if (pageCur.value >= pageTotal.value) {
            return
        }
        
        if (results.value[pageCur.value + 1]) {
            pageCur.value += 1
            return
        }
        
        
        try {
            loading.value = true
            const { data, err } = await searchTv(searchTerm, pageCur.value + 1)
            if (err) {
                throw err
            }
    
    
            pageTotal.value = data.TotalPages
            resultsCnt.value = data.TotalResults
            results.value = {
                ...results.value,
                [pageCur.value + 1]: data.Results
            }
            pageCur.value += 1
    
        } catch (error) {
            console.error(error)
        } finally {
            loading.value = false
        }
    
    }

    const fetchResults = async () => {
        // some logic
    }

    return {
        // data
        results,
        loading,
        pageCur,
        pageTotal,
        resultsCnt,

        // actions
        onSearch,
        onNxtBtn,
        onPrvBtn,

    }

})