<template>
    <v-row>
        <div>Directly</div>
    </v-row>
</template>

<script>

export default {
    name: "TabDirectly",
    components: {},
    props: {
      disabled: {
        Type: Boolean,
        required: true,
        default: false,
      },
      sessionNumber: {
        Type: Number,
        required: true,
        default: 0,
      },
      protocolObject: {
        Type: Object,
        required: true,
        default: 0,
      },
      typePanel: {
        Type: String,
        required: true,
        default: "",
      }
    }
}

</script>