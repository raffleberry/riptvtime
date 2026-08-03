import { useSearchStore } from "../tabs/Search/searchStore.js"
import { ref, storeToRefs } from "../vue.js"

const SearchButtons = {
  components: {},
  props: {},
  setup() {
    const store = useSearchStore()

    const { onSearch, onNxtBtn, onPrvBtn } = store

    const { loading, pageCur, pageTotal, resultsCnt } = storeToRefs(store)

    return {
      loading,
      pageCur,
      pageTotal,
      resultsCnt,

      onSearch,
      onNxtBtn,
      onPrvBtn,
    }
  },

  template: /* HTML */ `
    <div class="d-flex flex-column justify-content-center align-items-center">
      <div class="input-group mt-3 justify-content-center">
        <button
          :disabled="pageCur === 1 || loading"
          type="button"
          class="btn btn-outline-primary"
          @click="onPrvBtn"
        >
          <i class="bi bi-arrow-left"></i>
        </button>
        <div class="mx-3 d-flex flex-column align-items-center">
          <div>{{ pageTotal }} Results</div>
          <div>{{ pageCur }} / {{ pageTotal }}</div>
        </div>
        <button
          :disabled="pageCur === pageTotal || loading"
          type="button"
          class="btn btn-outline-primary"
          @click="onNxtBtn"
        >
          <i class="bi bi-arrow-right"></i>
        </button>
      </div>
    </div>
  `,
}

export { SearchButtons }
