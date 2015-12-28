package rethinkdb

import r "github.com/dancannon/gorethink"

// Insert uses the globally established session S to insert the `data` into the
// specified `table`. The `soft` parameter is used to run a Soft durability.
// This function returns the id of the new entry as the first returned element,
// or the error (second) in case of failure.
func Insert(table string, data interface{}, soft bool) (string, error) {
	opts := r.InsertOpts{}
	if soft {
		opts.Durability = "soft"
	} else {
		opts.Durability = "hard"
	}

	res, err := r.Table(table).Insert(data, opts).Run(S)
	if err != nil {
		return "", err
	}
	defer res.Close()

	var result RDict
	if err = res.One(&result); err != nil {
		return "", err
	}

	id := result["generated_keys"].([]interface{})[0].(string)
	return id, nil
}
