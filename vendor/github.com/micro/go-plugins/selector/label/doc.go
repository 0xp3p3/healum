// Package label is a priority label based selector.
package label

/*
   A priority based label selector. Rather than just returning nodes with specific labels
   this selector orders the nodes based on a list of labels. If no labels match all the
   nodes are still returned. The priority based label selector is useful for such things
   as rudimentary AZ based routing where requests made to other services should remain
   in the same AZ.
*/
