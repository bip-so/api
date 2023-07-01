package blocks

/*
We will add block type : simple_table
We can also add a flag has_children to block
We need a way to track columns on a given TABLE
We need limit the columns numbers and row numbers
	Table: (Container)
		- Column : 25
		- Rows : 200
	- Table Block : "simple_table"
	- Table Columns
		- We need to store Column Order.
		- We'll also store the BlocksID's
		- Row: [block1]
		- Columns [COL_ID1, COL_ID2, COL_ID3]

	- Table Rows: simple_table_row
		- Row 1:  COL_ID1 : ["apple"] ,  COL_ID2: ["orange"], COL_ID3 : ["banana"]
		- Row 2:  COL_ID4 : ["yellow"] ,  COL_ID5: ["orange"], COL_ID6 : ["dark yellow"]

Table

{
    type: 'table',
    children: [
      {
        type: 'table-row',
        children: [
          {
			coludId :
            type: 'table-cell',
            children: [{  }],
          },
          {
            type: 'table-cell',
            children: [{ text: 'Human', bold: true }],
          },
          {
            type: 'table-cell',
            children: [{ text: 'Dog', bold: true }],
          },
          {
            type: 'table-cell',
            children: [{ text: 'Cat', bold: true }],
          },
        ],
      },
      {
`
Rows





*/
