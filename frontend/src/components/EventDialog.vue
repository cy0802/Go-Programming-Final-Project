<template>
  <v-dialog v-model="dialog" max-width="500px">
    <v-card>
      <v-card-text>
        <v-text-field v-model="localEvent.title" label="事件名稱"></v-text-field>
        <v-row justify="center">
          <v-date-picker v-model="localEvent.start" v-if="newEvent" label="日期"></v-date-picker>
        </v-row>
        <v-textarea v-model="localEvent.diary" label="日記" class="mt-5"></v-textarea>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="save">Save</v-btn>
        <v-btn color="blue darken-1" text @click="close">Cancel</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, watch } from 'vue';
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  modelValue: Boolean,
  event: Object,
  newEvent: Boolean
});

const emit = defineEmits(['update:modelValue', 'save']);

const dialog = ref(props.modelValue);
const localEvent = ref({ ...props.event });
const newEvent = ref(false);

watch(() => props.modelValue, (newVal) => {
  dialog.value = newVal;
  localEvent.value = { ...props.event };
  if (localEvent.value.title){
    newEvent.value = false;
  } else {
    newEvent.value = true;
  }
});

// watch(() => localEvent.value, (newVal) => {
//   console.log(newVal);
// });

watch(dialog, (newVal) => {
  emit('update:modelValue', newVal);
});

const save = () => {
  emit('save', localEvent.value, newEvent.value);
  newEvent.value = false;
  close();
};

const close = () => {
  dialog.value = false;
  newEvent.value = false;
};
</script>