import { ky, PAGE, theme } from "../../utils.js"
import { computed, onMounted, ref } from "../../vue.js"

export const UpcomingTile = {
  props: {
    item: Object,
  },
  components: {},
  setup: (props) => {
    onMounted(() => {})

    const timeLeft = computed(() => {
      if (!props.item.Episode.AirDate) return null
      const date1 = new Date()
      const date2 = new Date(props.item.Episode.AirDate)
      const diffTime = Math.abs(date2 - date1)
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
      return diffDays
    })

    return {
      timeLeft,
      ky,
    }
  },
  template: /* HTML */ `
    <div class="card">
      <div class="row">
        <!--
        <div class="col-md-4">
        <img src="..." class="img-fluid rounded-start" alt="..."> 
        </div>
        <div class="col-md-8">
        </div>
        -->
        <div class="col">
          <div class="card-body">
            <h5 class="card-title">
              <router-link :to="'/series/' + item.Episode.SeriesMId">
                {{ item.SeriesName }} <span class="text-muted">({{ item.Year }})</span> - {{
                ky(item.Episode.Season, item.Episode.Episode) }}
              </router-link>
            </h5>
            <p class="card-text">{{ item.Episode.Overview }}</p>
            <p class="card-text text-start">
              {{ (new Date(item.Episode.AirDate)).toDateString() }}<span>
                (In {{ timeLeft }} days)</span
              >
            </p>
          </div>
        </div>
      </div>
    </div>
  `,
}
