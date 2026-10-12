import { apiGetUpcoming } from "../../api.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { PAGE, theme } from "../../utils.js"
import { onMounted, ref } from "../../vue.js"
import { UpcomingTile } from "./UpcomingTile.js"

const Upcoming = {
  props: {},
  components: {
    UpcomingTile,
  },
  setup: (props) => {
    const loading = ref(false)
    const upcoming = ref([])

    const fetchUpcoming = async () => {
      loading.value = true
      try {
        const { data, err } = await apiGetUpcoming()
        if (err) {
          throw err
        }
        upcoming.value = data || []
      } catch (err) {
        console.error(err)
        notify(MsgType.Error, "Upcoming", err)
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      fetchUpcoming()
    })

    return {
      loading,
      upcoming,
    }
  },
  template: /* HTML */ `
    <div>
      <h1 class="mb-4">Upcoming</h1>
      <div
        v-if="loading"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 50vh;"
      >
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div v-else>
        <div v-for="item in upcoming" class="mb-3">
          <UpcomingTile :item="item"></UpcomingTile>
        </div>
      </div>
    </div>
  `,
}
export { Upcoming }
