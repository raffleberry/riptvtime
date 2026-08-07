import { apiSearchTv, apiUploadImportZip } from "../../api.js"
import { notify } from "../../components/Notify/Notify.js"
import { PAGE, theme } from "../../utils.js"
import { computed, onMounted, ref, watch } from "../../vue.js"

const imgUrl = (id) => {
  if (!id) return ""
  return `https://image.tmdb.org/t/p/w342/${id}`
}

const Match = {
  props: {
    data: Object,
  },
  components: {},
  emits: ["matchDone"],
  setup: (props, ctx) => {
    onMounted(() => {})

    const searchTerm = ref("")
    const selected = ref(null)

    const tv = computed(() => {
      if (!props.series || props.series.length === 0) {
        return {}
      }
      return props.series[0]
    })

    const searchResults = ref([])
    const searchResultsImgLoaded = ref([])

    const loading = ref(false)

    const handleSearch = async () => {
      try {
        if (!searchTerm.value) return
        loading.value = true
        selected.value = null
        const { data, err } = await apiSearchTv(searchTerm.value, 1)
        if (err) {
          console.log(err)
          return
        }
        if (!data.Results) {
          searchResults.value = []
        } else {
          searchResults.value = data.Results
        }
        searchResultsImgLoaded.value = new Array(searchResults.value.length).fill(false)
      } catch (e) {
        notify(e)
      } finally {
        loading.value = false
      }
    }

    watch(
      () => props.series,
      () => {
        if (props.series.length === 0) return
        searchTerm.value = props.series[0].Name.replace(/ \(\s*\d{4}\s*\)\s*/g, "")
        handleSearch()
      },
      { immediate: true },
    )

    const onMatchDone = (TvTimeSId, MId) => {
      ctx.emit("matchDone", { TvTimeSId, MId })
    }

    return {
      tv,
      sr: searchResults,
      srl: searchResultsImgLoaded,
      imgUrl,
      loading,
      onMatchDone,
      selected,
      searchTerm,
      handleSearch,
    }
  },
  template: /* HTML */ `{data}`,
  templatex: /* HTML */ `
    <div>
      <h4 class="mt-4">Match ({{ series.length }} Pending) - {{ tv.Name }}</h4>
      <div class="col">
        <div class="d-flex flex-row align-items-center my-2">
          <button
            v-if="selected"
            class="btn btn-success me-2"
            @click="onMatchDone(tv.TvTimeSId, selected.Id)"
          >
            Confirm
          </button>
          <span class="align-middle me-2" v-if="selected"
            >Selected: {{selected.Name}} ({{selected.Year}})</span
          >

          <p v-if="!selected">Please Select the correct one:</p>
        </div>

        <div class="row my-2 px-0">
          <div class="input-group search-container">
            <input
              v-model="searchTerm"
              type="text"
              class="form-control search-input"
              placeholder="Search..."
              @keyup.enter="handleSearch"
            />
            <button
              class="btn btn-outline-primary"
              type="button"
              id="searchButton"
              @click="handleSearch"
            >
              <i class="bi bi-search"></i>
            </button>
          </div>
        </div>
        <div
          v-if="loading"
          class="d-flex flex-row justify-content-center align-items-center"
          style="min-height: 320px;"
        >
          <div class="spinner-border text-secondary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <div
          v-if="!loading && sr.length == 0"
          class="d-flex justify-content-center align-items-center"
        >
          <h2>No Results</h2>
        </div>
        <div
          v-if="!loading && sr.length > 0"
          class="mx-0 pb-2 row d-flex flex-nowrap overflow-x-auto"
        >
          <div
            v-for="(r, idx) in sr"
            class="me-3 d-flex flex-column justify-content-between card p-1"
            :class="{ 'bg-primary': r.Id === selected?.Id }"
            style="width: 240px; cursor: pointer;"
            @click="selected = r"
          >
            <div class="position-relative" style="min-height: 320px;">
              <div
                v-if="!srl[idx] && r.Image"
                class="position-absolute top-50 start-50 translate-middle spinner-border text-secondary"
                role="status"
              >
                <span class="visually-hidden">Loading...</span>
              </div>
              <img
                :src="imgUrl(r.Image)"
                loading="lazy"
                class="img-fluid rounded"
                :alt="r.Name"
                @load="srl[idx] = true"
              />
            </div>
            <div class="text-center mt-2">{{ r.Name }} ({{r.Year}})</div>
          </div>
        </div>
        <div v-if="selected" class="row"><p>Overview: {{selected?.Overview}}</p></div>
      </div>
    </div>
  `,
}

export { Match }
