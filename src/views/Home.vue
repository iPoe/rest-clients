<template>
  <div class="dashboard">
    <h1>Sin miedo al exito</h1>
    <v-container class="grey lighten-5">
      <v-row justify="start">
        <v-col md ="4" >
          <Date v-on:date="getDate($event)" />
       <v-btn
        color="primary"
        elevation="2"
        @click="printdate()"
      >click me
      <v-icon right>search</v-icon>
      </v-btn>

        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col md ="4" >
        <v-card
          class="pa-2"
          outlined
          tile
        >
          Lista de clientes
        </v-card>

        </v-col>
      </v-row>
      

    </v-container>
  </div>
</template>

<script>

import Date from '@/components/Date.vue'

export default {
  data(){
    return{
      date:null
    }
  },
  components:{
    Date
  },
  methods:{
    getDate(e) {
      this.date = e
    },
    printdate(){
      console.log(this.date)
    }
  }
      
}
</script>
