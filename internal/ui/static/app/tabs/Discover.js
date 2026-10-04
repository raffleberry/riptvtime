import { computed, onMounted, storeToRefs, useRoute } from "../vue.js"
import { DiscoverRail } from "./DiscoverRail.js"
import { DiscoverViewAll } from "./DiscoverViewAll.js"
import { DISCOVER_LISTS, useDiscoverStore } from "./discoverStore.js"

const Discover = {
  props: {},
  components: {
    DiscoverRail,
    DiscoverViewAll,
  },
  setup: (props) => {
    const store = useDiscoverStore()
    const { rails } = storeToRefs(store)
    const { fetchRails } = store

    const route = useRoute()
    const view = computed(() => route.query.view || "")

    onMounted(() => {
      fetchRails()
    })

    const railState = (key) => {
      return (
        rails.value[key] || { items: [], totalPages: 1, totalResults: 0, loading: true, err: null }
      )
    }

    return {
      DISCOVER_LISTS,
      view,
      railState,
    }
  },
  template: /* HTML */ `
    <div>
      <DiscoverViewAll v-if="view"></DiscoverViewAll>
      <div v-else>
        <h1 class="px-3 my-2">Discover</h1>
        <DiscoverRail
          v-for="l in DISCOVER_LISTS"
          :key="l.key"
          :title="l.title"
          :listKey="l.key"
          :items="railState(l.key).items"
          :loading="railState(l.key).loading"
        >
        </DiscoverRail>
      </div>
    </div>
  `,
}
export { Discover }
