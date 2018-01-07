type 'a invalid = {
  given : 'a;
  valid : 'a;
}

exception Invalid_arity of int invalid

exception Invalid_type of string invalid

let invalid_arity ~given ~valid =
  Invalid_arity { given; valid }

let invalid_type ~given ~valid =
  Invalid_type { given; valid }

