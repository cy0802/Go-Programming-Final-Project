<template>
  <v-app>
    <v-toolbar 
      title="時光回憶錄"
      color="primary"
    >
      <v-btn
        class="mr-5"
        @click="handleDump()"
        variant="tonal"
      >
        DUMP
      </v-btn>
    </v-toolbar>
    <v-row>
      <v-col cols="1"></v-col>
      <v-col cols="10">
        <v-sheet height="600">
          <v-calendar
            :events="events"
            hide-week-number
            next-icon="mdi-chevron-right"
            prev-icon="mdi-chevron-left"
            class="my-1"
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
      <v-col cols="1"></v-col>
      <EventDialog
        v-model="dialog"
        :event="selectedEvent"
        @save="updateEvent"
        @deleteEvent="deleteEvent"
      />
      <DumpDialog
        v-model="dumpDialog"
        :dumpData="dumpData"
        :loading="loading"
        @close="dumpDialog = false"
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
  </v-app>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useDate } from 'vuetify'
import EventDialog from './components/EventDialog.vue';
import DumpDialog from './components/DumpDialog.vue';
import * as eventService from './services/event.js'
import { getDump } from './services/dump.js'

const adapter = useDate();
const events = ref([]);
const selectedEvent = ref({});
const dialog = ref(false);
const dumpDialog = ref(false);
const loading = ref(false);
const dumpData = ref('');

async function handleDump() {
  try {
    dumpDialog.value = true;
    loading.value = true;
    dumpData.value = await getDump();
    loading.value = false;
  } catch (error) {
    console.error("Error dumping data:", error);
  }
}

function addEvent(id, name, start, end, diary = '') {
  events.value.push({ 
    id: id,
    title: name,
    start: start,
    end: end,
    color: 'cyan',
    allDay: 1,
    diary: diary
  });
}

function deleteEvent(id) {
  try {
    eventService.deleteEvent(id);
    const index = events.value.findIndex(event => event.id === id);
    events.value.splice(index, 1);
  } catch (error) {
    console.error("Error deleting event:", error);
  }
}

function openDialog(event) {
  selectedEvent.value = { ...event };
  dialog.value = true;
}

function updateEvent(updatedEvent, newEvent) {
  const formattedDate = new Date(updatedEvent.start).toLocaleDateString('zh-TW').replaceAll('/', '-');
  if (newEvent) {
    try {
      eventService.createEvent({ title: updatedEvent.title, date: formattedDate, diary: updatedEvent.diary });
      addEvent(events.value.length + 1, updatedEvent.title, updatedEvent.start, updatedEvent.start, updatedEvent.diary);
    } catch (error) {
      console.error("Error creating event:", error);
    }
  } else {
    try {
      eventService.updateEvent(updatedEvent);
      const index = events.value.findIndex(event => event.id === selectedEvent.value.id);
      events.value[index].title = updatedEvent.title;
      events.value[index].diary = updatedEvent.diary;
    } catch (error) {
      console.error("Error updating event:", error);
    }
  }
}

function addNewEvent() {
  selectedEvent.value = {};
  dialog.value = true;
}

onMounted(async () => {
  try {
    var fetchedEvent = await eventService.getEvents();
    fetchedEvent.forEach(event => {
      addEvent(event.id, event.title, adapter.startOfDay(new Date(event.date)), adapter.endOfDay(new Date(event.date)), event.diary);
    });
  } catch (error) {
    console.error("Error fetching events:", error);
  }
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