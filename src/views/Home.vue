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
        @click="fetchClientsbyDate(date)"
      >Search
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
import axios from 'axios'
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
    },
    fetchClientsbyDate(date){
      const [year, month, day] = date.split('-')
      console.log(this.date)
			axios.get("http://localhost:9000/data/?date="+`${month}/${day}/${year}`).then((res)=>{
				// this.datos = res.data["datos"];
        console.log(res.data)
			}).catch((error) =>{
				this.$vs.notify({
					color:'danger',
					title:'Error updating db',
					text: error,
					iconPack: 'feather', icon:'icon-alert-circle'
				});
				console.log(error);
			});
		}
  }
      
}
</script>
