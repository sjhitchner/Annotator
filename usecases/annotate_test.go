package usecases

import (
	"fmt"
	"github.com/sjhitchner/sourcegraph/infrastructure"
	"testing"
)

/*
echo >> $OUT
echo '=== Test simple annotation ===' >> $OUT
curl -XDELETE "http://$HOST:$PORT/names"
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/alex" -d '{ "url": "http://alex.com" }'
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d  | formatHTML >> $OUT
*/
func TestSimpleAnnotation(t *testing.T) {
	repo := infrastructure.NewNameRepository()
	interactor := NewAnnotationInteractor(repo)

	if err := repo.Put("alex", "http://alex.com"); err != nil {
		t.Fatal(err)
	}

	input = `my name is alex`
	expected = `my name is <a href=http://alex.com>alex</a>`

	html, err := interactor.AnnotateHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if html != expected {
		t.Fatalf("expected [%s] != result [%s]", expected, html)
	}
}

/*
echo >> $OUT
echo '=== Test annotation of multiple names ===' >> $OUT
curl -XDELETE "http://$HOST:$PORT/names"
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/alex" -d '{ "url": "http://alex.com" }'
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/bo" -d '{ "url": "http://bo.com" }'
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/casey" -d '{ "url": "http://casey.com" }'
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d 'alex, bo, and casey went to the park.' | formatHTML >> $OUT
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d 'alex alexander alexandria alexbocasey' | formatHTML >> $OUT

echo >> $OUT
echo '=== Test HTML correctness ===' >> $OUT
curl -XDELETE "http://$HOST:$PORT/names"
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/alex" -d '{ "url": "http://alex.com" }'
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d '<div data-alex="alex">alex</div>' | formatHTML >> $OUT
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d '<a href="http://foo.com">alex is already linked</a> but alex is not' | formatHTML >> $OUT
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d "<div><p>this is paragraph 1 about alex.</p><p>alex's paragraph number 2.</p><p>and some closing remarks about alex</p></div>" | formatHTML >> $OUT

echo >> $OUT
echo '=== Test additional annotations ===' >> $OUT
curl -XDELETE "http://$HOST:$PORT/names"
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/alex" -d '{ "url": "http://alex.com" }'
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/bo" -d '{ "url": "http://bo.com" }'
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/casey" -d '{ "url": "http://casey.com" }'
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d '<div data-alex="alex">alex</div>' | formatHTML >> $OUT
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d "<div><p>this is paragraph 1 about alex.</p><p>alex's paragraph number 2.</p><p>and some closing remarks about alex</p></div>" | formatHTML >> $OUT
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d "<div><ul><li>alex</li><li>bo</li><li>bob</li><li>casey</li></ul></div><div><p>this is paragraph 1 about alex.</p><p>alex's paragraph number 2.</p><p>and some closing remarks about alex</p></div>" | formatHTML >> $OUT

echo >> $OUT
echo '=== Test annotation of complex example ===' >> $OUT
curl -XDELETE "http://$HOST:$PORT/names"
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/Sourcegraph" -d '{ "url": "https://sourcegraph.com" }'
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/Milton" -d '{ "url": "https://www.google.com/search?q=milton" }'
curl -H 'Content-Type:application/json' -XPUT "http://$HOST:$PORT/names/strong" -d '{ "url": "https://www.google.com/search?q=strong" }'
curl -H 'Content-Type:text/plain' -XPOST "http://$HOST:$PORT/annotate" -d "<div class=\"row\"><div class=\"col-md-6\"><h2> Sourcegraph makes programming <strong>delightful.</strong></h2><p>We want to make you even better at what you do best: building software to solve real problems.</p><p>Sourcegraph makes it easier to find the information you need: documentation, examples, usage statistics, answers, and more.</p><p>We're just getting started, and we'd love to hear from you. <a ui-sref=\"help.contact\" href=\"/contact\">Get in touch with us.</a></p></div><div class=\"col-md-4 team\"><h3>Team</h3><ul><li><img src=\"https://secure.gravatar.com/avatar/c728a3085fc16da7c594903ea8e8858f?s=64\" class=\"pull-left\"><div class=\"bio\"><strong>Beyang Liu</strong><br><a target=\"_blank\" href=\"http://github.com/beyang\">github.com/beyang</a><a href=\"mailto:beyang@sourcegraph.com\">beyang@sourcegraph.com</a></div></li><li><img src=\"https://secure.gravatar.com/avatar/d491971c742b8249341e495cf53045ea?s=64\" class=\"pull-left\"><div class=\"bio\"><strong>Quinn Slack</strong><br><a target=\"_blank\" href=\"http://github.com/sqs\">github.com/sqs</a><a href=\"mailto:sqs@sourcegraph.com\">sqs@sourcegraph.com</a></div></li><li><img src=\"https://1.gravatar.com/avatar/43ec631d6fda6a1cf42aaf875d784597?d=https%3A%2F%2Fidenticons.github.com%2F71945c68441f29a222b5689f640c956f.png&amp;r=x&amp;s=440\" class=\"pull-left\"><div class=\"bio\"><strong>Yin Wang</strong><br><a target=\"_blank\" href=\"http://github.com/yinwang0\">github.com/yinwang0</a><a target=\"_blank\" href=\"http://yinwang0.wordpress.com\">yinwang0.wordpress.com</a><a href=\"mailto:yin@sourcegraph.com\">yin@sourcegraph.com</a></div></li><li><img src=\"https://s3-us-west-2.amazonaws.com/public-dev/milton.png\" class=\"pull-left\"><div class=\"bio\"><strong>Milton</strong> the Australian Shepherd </div></li></ul><p><a ui-sref=\"help.contact\" href=\"/contact\">Want to join us?</a></p></div></div>" | formatHTML >> $OUT
*/
