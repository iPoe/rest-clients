<template>
    <v-row>
    <v-menu
        ref="menu"
        v-model="menu"
        :close-on-content-click="false"
        :return-value.sync="date"
        transition="scale-transition"
        offset-y
        min-width="auto"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
            v-model="date"
            label="Choose a date"
            prepend-icon="mdi-calendar"
            readonly
            v-bind="attrs"
            v-on="on"
          ></v-text-field>
        </template>
        <v-date-picker
          v-model="date"
          hint="DD/MM/YYYY format"
          no-title
          scrollable
          type="date"
        >
          <v-spacer></v-spacer>
          <v-btn
            text
            color="primary"
            @click="menu = false"
          >
            Cancel
          </v-btn>
          <v-btn
            text
            color="primary"
            @click="$refs.menu.save(date);master()"
          >
            OK
          </v-btn>
        </v-date-picker>
      </v-menu>
      

    </v-row>
</template>

<script>
export default {
    data: () =>  ({
      date: new Date().toISOString().substr(0, 10),
      menu: false,
      datos: null,
    }),
    methods:{
      master(){
        
        this.$emit('date', this.date)        
      },   
    },
    
   
  }
</script>
