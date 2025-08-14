<script setup>
import { ref, computed, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import useProjectAPI from '../../api/Project';

const route = useRoute()
const projectAPI = useProjectAPI();

const eventAPI = () => projectAPI.eventAPIFor(route.params.project_id);

const occurrences = ref([])
const exceptions = ref([])

const startDate = new Date();
const endDate = new Date(new Date().setFullYear(new Date().getFullYear() + 1));

watchEffect(async () => {
  occurrences.value = await eventAPI().getOccurrences(route.params.id, startDate, endDate);
  exceptions.value = await eventAPI().getExceptions(route.params.id, startDate, endDate);
});

const reload = () => {
  eventAPI().getOccurrences(route.params.id, startDate, endDate).then((data) => {
    occurrences.value = data;
  });
  eventAPI().getExceptions(route.params.id, startDate, endDate).then((data) => {
    exceptions.value = data;
  });
};

const removeException = async (exception_id) => {
  await eventAPI().deleteException(route.params.id, exception_id);
  reload()
};

const addException = async (occurrence_id) => {
  await eventAPI().createException(route.params.id, occurrence_id);
  reload()
};

const ocurrencesByMonthYear = computed(() => {

  const occurrencesMap = {};
  occurrences.value.forEach(occurrence => {
    const occurrenceDate = new Date(occurrence.start_date);
    const monthYearKey = `${occurrenceDate.getMonth() + 1}-${occurrenceDate.getFullYear()}`;
    if (!occurrencesMap[monthYearKey]) {
      occurrencesMap[monthYearKey] = [];
    }
    occurrencesMap[monthYearKey].push({ id: "occurrence-" + occurrence.id, occurrence_id: occurrence.id, start_date: occurrence.start_date, end_date: occurrence.end_date, isException: false });
  });

  exceptions.value.forEach(exception => {
    const monthYearKey = `${new Date(exception.start_date).getMonth() + 1}-${new Date(exception.start_date).getFullYear()}`;
    if (!occurrencesMap[monthYearKey]) {
      occurrencesMap[monthYearKey] = [];
    }
    occurrencesMap[monthYearKey].push({ id: "exception-" + exception.id, exception_id: exception.id, start_date: exception.start_date, end_date: exception.end_date, isException: true });
  });

  Object.keys(occurrencesMap).forEach(monthYearKey => {
    occurrencesMap[monthYearKey] = occurrencesMap[monthYearKey].sort((a, b) => new Date(a.start_date) - new Date(b.start_date));
  });

  return occurrencesMap;
});

const getGermanMonthName = (monthNumber) => {
  const monthNames = [
    "Januar", "Februar", "März", "April", "Mai", "Juni",
    "Juli", "August", "September", "Oktober", "November", "Dezember"
  ];
  return monthNames[monthNumber - 1];
};

const displayMonthYear = (monthYear) => {
  const [month, year] = monthYear.split('-');
  return `${getGermanMonthName(parseInt(month))} ${year}`;
};

const formatTimespan = (dates) => {
  if (dates.length == 0) {
    return ""
  }
  let startDate = new Date(dates[0]);
  let endDate = new Date(dates[1]);
  let formattedStartDate = startDate.getDate() + "." + (startDate.getMonth() + 1) + "." + startDate.getFullYear() + " " + startDate.getHours().toString().padStart(2, '0') + ":" + startDate.getMinutes().toString().padStart(2, '0');
  let formattedEndDate = endDate.getDate() + "." + (endDate.getMonth() + 1) + "." + endDate.getFullYear() + " " + endDate.getHours().toString().padStart(2, '0') + ":" + endDate.getMinutes().toString().padStart(2, '0');
  return formattedStartDate + " - " + formattedEndDate;
}

</script>

<template>
  <h1>Termine &amp; Ausnahmen</h1>
  <div class="container">
    <div v-for="monthYear in Object.keys(ocurrencesByMonthYear).sort((a, b) => new Date(a.split('-')[1], a.split('-')[0]) - new Date(b.split('-')[1], b.split('-')[0]))"
      :key="monthYear">
      <div class="row">
        <div class="col">
          <h3>{{ displayMonthYear(monthYear) }}</h3>
        </div>
      </div>
      <div class="row" v-for="occurrence in ocurrencesByMonthYear[monthYear]" :key="occurrence.id">
        <div class="col-md-6">
          <p v-if="occurrence.isException"><s>{{ formatTimespan([occurrence.start_date, occurrence.end_date]) }}</s></p>
          <p v-else>{{ formatTimespan([occurrence.start_date, occurrence.end_date]) }}</p>
        </div>
        <div class="col-md-6">
          <button v-if="occurrence.isException" type="button" class="btn btn-success"
            @click="removeException(occurrence.exception_id)">Wiederherstellen</button>
          <button v-else type="button" class="btn btn-danger" @click="addException(occurrence.occurrence_id)">Ausnahme
            hinzufügen</button>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <div class="functions" style="margin-top: 2rem">
          <button type="button" class="btn btn-secondary" @click="$router.go(-1)">Zurück</button>
        </div>
      </div>
    </div>
  </div>
</template>