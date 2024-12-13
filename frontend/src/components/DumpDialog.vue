<template>
  <v-dialog 
    v-model="localDialog" 
    max-width="500px"
    persistent
  >
    <v-card>
      <v-card-title
        class="mt-2"
      >
        本周總結
      </v-card-title>
      <v-card-text class="mt-5">
        <v-row justify="center">
          <v-progress-circular
            v-if="localLoading"
            indeterminate
            color="primary"
          />
        </v-row>
        <v-img 
          v-if="!localLoading"
          :src="localDumpData"
          cover
        />
      </v-card-text>
      <v-card-actions>
        <v-btn @click="close()">關閉</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, defineEmits, defineProps, watch } from 'vue';

const props = defineProps({
  dialog: Boolean,
  dumpData: String,
  loading: Boolean
});

const emit = defineEmits(['close']);  

const localDialog = ref(props.dialog);
const localDumpData = ref(props.dumpData);
const localLoading = ref(props.loading);

watch(() => props.dialog, (newVal) => {
  localDialog.value = newVal;
});

watch(() => props.dumpData, (newVal) => {
  localDumpData.value = newVal;
});

watch(() => props.loading, (newVal) => {
  localLoading.value = newVal;
});

const close = () => {
  emit('close');
};

</script>