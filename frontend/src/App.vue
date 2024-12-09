<template>
  <v-row class="fill-height">
    <v-col>
      <v-sheet height="600">
        <v-calendar
          :events="events"
          hide-week-number
          next-icon="mdi-chevron-right"
          prev-icon="mdi-chevron-left"
        >
          <template #event="{ event }">
            <div class="center-content">
              <v-chip 
                @click="openDialog(event)" 
                class="event chip-style my-1"
                color="primary"
                variant="elevated"
                label
              >
                {{ event.title }}
              </v-chip>
            </div>
          </template>
        </v-calendar>
      </v-sheet>
    </v-col>
    <EventDialog
      v-model="dialog"
      :event="selectedEvent"
      @save="updateEvent"
    />
    <v-btn
      location="bottom right"
      position="fixed"
      class="mx-10 my-5"
      size="60"
      color="primary"
      icon="mdi-plus"
      floating
      @click="addNewEvent()"
    />
  </v-row>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useDate } from 'vuetify'
import EventDialog from './components/EventDialog.vue';
const adapter = useDate();
const events = ref([]);
const selectedEvent = ref({});
const dialog = ref(false);

function addEvent(name, start, end, diary = '') {
  events.value.push({ 
    title: name,
    start: start,
    end: end,
    color: 'cyan',
    allDay: 1,
    diary: diary
  });
}

function openDialog(event) {
  selectedEvent.value = { ...event };
  dialog.value = true;
}

function updateEvent(updatedEvent, newEvent) {
  console.log("updateEvent called");
  console.log(updatedEvent);
  console.log(newEvent);
}

function addNewEvent() {
  selectedEvent.value = {};
  dialog.value = true;
}

onMounted(() => {
  addEvent('Event 2', adapter.startOfDay(new Date('2024-12-01')), adapter.endOfDay(new Date('2024-12-01')), "hello, world");
  addEvent('Event 4', adapter.startOfDay(new Date('2024-12-01')), adapter.endOfDay(new Date('2024-12-01')));
  addEvent('Event 2024', adapter.startOfDay(new Date('2024-12-02')), adapter.endOfDay(new Date('2024-12-02')));
});

</script>

<style>
.chip-style {
  width: 90%;
}
.center-content {
  display: flex;
  justify-content: center;
}
</style>