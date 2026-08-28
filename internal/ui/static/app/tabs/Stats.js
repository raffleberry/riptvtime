import { apiGetSeriesStats } from "../api.js"
import { notify } from "../components/Notify/Notify.js"
import { PAGE, theme } from "../utils.js"
import { onMounted, ref } from "../vue.js"

const Stats = {
  props: {},
  components: {},
  setup: (props) => {
    const stats = ref({})

    const fetchStats = async () => {
      const { data, err } = await apiGetSeriesStats()
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Stats", err)
        return
      }
      stats.value = data
    }

    onMounted(() => {
      fetchStats()
    })

    return {
      stats,
    }
  },
  template: /* HTML */ `
    <div>
      <div class="d-flex align-items-center">
        <div class="card text-center me-2">
          <div class="card-body">
            <h5 class="card-title">{{ Math.ceil(stats.TotalHours) }}</h5>
            <p class="card-text">Hours Watched</p>
          </div>
        </div>
        <div class="card text-center me-2">
          <div class="card-body">
            <h5 class="card-title">{{ Math.ceil(stats.TotalEpisodes) }}</h5>
            <p class="card-text">Episodes Watched</p>
          </div>
        </div>
        <div class="card text-center me-2">
          <div class="card-body">
            <h5 class="card-title">{{ Math.ceil(stats.TotalShows) }}</h5>
            <p class="card-text">Shows Tracked</p>
          </div>
        </div>
      </div>
      <div></div>
    </div>
  `,
}
export { Stats }
