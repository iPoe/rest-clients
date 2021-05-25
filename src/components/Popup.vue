<template>
<div>
<v-card hover>
    <v-container>
          <v-card >
            <v-row>
            <v-col cols=7>
            <v-toolbar
              color="teal"
              dark
            >
              <v-app-bar-nav-icon></v-app-bar-nav-icon>

              <v-toolbar-title>User Transactions</v-toolbar-title>

              <v-spacer></v-spacer>

              <v-btn icon>
                <v-icon>dots-vertical</v-icon>
              </v-btn>
            </v-toolbar>
            <v-list
            style="max-height: 250px"
            class="overflow-y-auto"
            >
      <v-list-group
        v-for="item in Tdata"
        :key="item.Tid"
        no-action
      >
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title >ID {{item.Tid}}</v-list-item-title>
              </v-list-item-content>
                <v-spacer />
              <v-list-item-content>
                <v-list-item-title >TOTAL {{ item.price | toCurrency }}</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item
              v-for="child in item.ProductIds"
              :key="child"
            >
              <v-list-item-content>
                <v-list-item-title v-text="child"></v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
        </v-list>
        </v-col>
         <v-col cols=4>
      <v-list style="max-height: 150px"
            class="overflow-y-auto">
        <v-subheader>Users with same IP</v-subheader>
        <v-list-item-group
          color="primary"
        >
          <v-list-item
            v-for="(item, i) in similarBuyer"
            :key="i"
          >
            <v-list-item-icon>
              <v-icon >face</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title v-text="item"></v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
      </v-col>
      </v-row>
      </v-card>
      

    </v-container>

      <v-card-actions>  
        <v-icon>email</v-icon>
        <v-spacer />
        <v-btn color="blue lighten-3">I'm fake</v-btn>
      </v-card-actions>
      
    </v-card>
   

</div>  
</template>

<script>
//Montar el bus y recibir la info
export default {
  props:['Tdata','similarBuyer'],
    data () {
      return {
      Theader:[
        { text: 'Transaction No', value: 'Tid' },
        { text: 'Products', value: 'ProductIds' },
        { text: 'TOTAL', value: 'price' },
      ]
      }
    },
    
}
</script>