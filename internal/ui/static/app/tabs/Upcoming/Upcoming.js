import { apiGetUpcoming } from "../../api.js"
import { Empty } from "../../components/Empty.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { onMounted, ref } from "../../vue.js"
import { UpcomingTile } from "./UpcomingTile.js"

const Upcoming = {
  props: {},
  components: {
    UpcomingTile,
    Empty
  },
  setup: (props) => {
    const loading = ref(false)
    const upcoming = ref([])

    const fetchUpcoming = async () => {
      loading.value = true

      const { data, err } = await apiGetUpcoming()
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Upcoming", err)
        return
      }
      upcoming.value = data

      loading.value = false
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
        <div v-if="upcoming.length === 0">
          <Empty icon="bi bi-send-check-fill"
            title="All caught up!"
            message="Explore the Discover Page to find more shows you'll love."
          />
        </div>
        <div v-else v-for="item in upcoming" class="mb-3">
          <UpcomingTile :item="item"></UpcomingTile>
        </div>
      </div>
    </div>
  `,
}
export { Upcoming }
