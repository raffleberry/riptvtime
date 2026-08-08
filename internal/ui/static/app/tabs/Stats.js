import { apiGetStatsTotal } from "../api.js"
import { notify } from "../components/Notify/Notify.js"
import { PAGE, theme } from "../utils.js"
import { onMounted, ref } from "../vue.js"

const Stats = {
  props: {},
  components: {},
  setup: (props) => {
    const total = ref({})

    const fetchStats = async () => {
      const { data, err } = await apiGetStatsTotal()
      if (err) {
        console.log(err)
        notify(MsgType.Error, "Stats", err)
        return
      }
      total.value = data
    }

    onMounted(() => {
      fetchStats()
    })

    return {
      total,
    }
  },
  template: /* HTML */ `
    <div>
      <h1>Stats</h1>
      <h4>Total Hours Watched: {{ (total.Total / 60).toFixed(2) }} Hours</h4>
    </div>
  `,
}
export { Stats }
