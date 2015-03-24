package usecases

import (
	. "github.com/sjhitchner/annotator/domain"
	"github.com/sjhitchner/annotator/interfaces/db"
	"testing"
)

func setupInteractor(params ...string) AnnotationInteractor {
	repo := db.NewNamesRepository()

	if len(params)%2 != 0 {
		panic("invalid params")
	}

	for i := 0; i < len(params); i += 2 {
		if err := repo.Put(Name(params[i]), URL(params[i+1])); err != nil {
			panic(err)
		}
	}

	return NewAnnotationInteractor(repo)
}

func annotateTest(t *testing.T, interactor AnnotationInteractor, input, expected string) {
	html, err := interactor.AnnotateHTML(input)
	if err != nil {
		t.Fatal(err)
	}
	if html != expected {
		t.Fatalf("\nexpected:\n[%s]\nresult:\n[%s]", expected, html)
	}
}

func TestSimpleAnnotation(t *testing.T) {
	interactor := setupInteractor(
		"alex", "http://alex.com",
	)

	annotateTest(
		t,
		interactor,
		`my name is alex`,
		`my name is <a href="http://alex.com">alex</a>`,
	)
}

func TestAnnotationOfMultipleNames(t *testing.T) {
	interactor := setupInteractor(
		"alex", "http://alex.com",
		"bo", "http://bo.com",
		"casey", "http://casey.com",
	)

	annotateTest(
		t,
		interactor,
		`alex, bo, and casey went to the park.`,
		`<a href="http://alex.com">alex</a>, <a href="http://bo.com">bo</a>, and <a href="http://casey.com">casey</a> went to the park.`,
	)
	annotateTest(
		t,
		interactor,
		`alex alexander alexandria alexbocasey`,
		`<a href="http://alex.com">alex</a> alexander alexandria alexbocasey`,
	)
}

func TestHTMLCorrectness(t *testing.T) {
	interactor := setupInteractor(
		"alex", "http://alex.com",
	)

	annotateTest(
		t,
		interactor,
		`<div data-alex="alex">alex</div>`,
		`<div data-alex="alex"><a href="http://alex.com">alex</a></div>`,
	)
	annotateTest(
		t,
		interactor,
		`<a href="http://foo.com">alex is already linked</a> but alex is not`,
		`<a href="http://foo.com">alex is already linked</a> but <a href="http://alex.com">alex</a> is not`,
	)
	annotateTest(
		t,
		interactor,
		`<div><p>this is paragraph 1 about alex.</p><p>alex's paragraph number 2.</p><p>and some closing remarks about alex</p></div>`,
		`<div><p>this is paragraph 1 about <a href="http://alex.com">alex</a>.</p><p><a href="http://alex.com">alex</a>'s paragraph number 2.</p><p>and some closing remarks about <a href="http://alex.com">alex</a></p></div>`,
	)
}

func TestAdditionalAnnotations(t *testing.T) {
	interactor := setupInteractor(
		"alex", "http://alex.com",
		"bo", "http://bo.com",
		"casey", "http://casey.com",
	)

	annotateTest(
		t,
		interactor,
		`<div data-alex="alex">alex</div>`,
		`<div data-alex="alex"><a href="http://alex.com">alex</a></div>`,
	)
	annotateTest(
		t,
		interactor,
		`<div><p>this is paragraph 1 about alex.</p><p>alex's paragraph number 2.</p><p>and some closing remarks about alex</p></div>`,
		`<div><p>this is paragraph 1 about <a href="http://alex.com">alex</a>.</p><p><a href="http://alex.com">alex</a>'s paragraph number 2.</p><p>and some closing remarks about <a href="http://alex.com">alex</a></p></div>`,
	)
	annotateTest(
		t,
		interactor,
		`<div><ul><li>alex</li><li>bo</li><li>bob</li><li>casey</li></ul></div><div><p>this is paragraph 1 about alex.</p><p>alex's paragraph number 2.</p><p>and some closing remarks about alex</p></div>`,
		`<div><ul><li><a href="http://alex.com">alex</a></li><li><a href="http://bo.com">bo</a></li><li>bob</li><li><a href="http://casey.com">casey</a></li></ul></div><div><p>this is paragraph 1 about <a href="http://alex.com">alex</a>.</p><p><a href="http://alex.com">alex</a>'s paragraph number 2.</p><p>and some closing remarks about <a href="http://alex.com">alex</a></p></div>`,
	)
}

func TestAnnotationOfComplexExample(t *testing.T) {
	interactor := setupInteractor(
		"Sourcegraph", "https://sourcegraph.com",
		"Milton", "https://www.google.com/search?q=milton",
		"strong", "https://www.google.com/search?q=strong",
	)

	annotateTest(
		t,
		interactor,
		`<div class="row"><div class="col-md-6"><h2> Sourcegraph makes programming <strong>delightful.</strong></h2><p>We want to make you even better at what you do best: building software to solve real problems.</p><p>Sourcegraph makes it easier to find the information you need: documentation, examples, usage statistics, answers, and more.</p><p>We're just getting started, and we'd love to hear from you. <a ui-sref="help.contact" href="/contact">Get in touch with us.</a></p></div><div class="col-md-4 team"><h3>Team</h3><ul><li><img src="https://secure.gravatar.com/avatar/c728a3085fc16da7c594903ea8e8858f?s=64" class="pull-left"><div class="bio"><strong>Beyang Liu</strong><br><a target="_blank" href="http://github.com/beyang">github.com/beyang</a><a href="mailto:beyang@sourcegraph.com">beyang@sourcegraph.com</a></div></li><li><img src="https://secure.gravatar.com/avatar/d491971c742b8249341e495cf53045ea?s=64" class="pull-left"><div class="bio"><strong>Quinn Slack</strong><br><a target="_blank" href="http://github.com/sqs">github.com/sqs</a><a href="mailto:sqs@sourcegraph.com">sqs@sourcegraph.com</a></div></li><li><img src="https://1.gravatar.com/avatar/43ec631d6fda6a1cf42aaf875d784597?d=https%3A%2F%2Fidenticons.github.com%2F71945c68441f29a222b5689f640c956f.png&amp;r=x&amp;s=440" class="pull-left"><div class="bio"><strong>Yin Wang</strong><br><a target="_blank" href="http://github.com/yinwang0">github.com/yinwang0</a><a target="_blank" href="http://yinwang0.wordpress.com">yinwang0.wordpress.com</a><a href="mailto:yin@sourcegraph.com">yin@sourcegraph.com</a></div></li><li><img src="https://s3-us-west-2.amazonaws.com/public-dev/milton.png" class="pull-left"><div class="bio"><strong>Milton</strong> the Australian Shepherd </div></li></ul><p><a ui-sref="help.contact" href="/contact">Want to join us?</a></p></div></div>`,
		`<div class="row"><div class="col-md-6"><h2> <a href="https://sourcegraph.com">Sourcegraph</a> makes programming <strong>delightful.</strong></h2><p>We want to make you even better at what you do best: building software to solve real problems.</p><p><a href="https://sourcegraph.com">Sourcegraph</a> makes it easier to find the information you need: documentation, examples, usage statistics, answers, and more.</p><p>We're just getting started, and we'd love to hear from you. <a ui-sref="help.contact" href="/contact">Get in touch with us.</a></p></div><div class="col-md-4 team"><h3>Team</h3><ul><li><img src="https://secure.gravatar.com/avatar/c728a3085fc16da7c594903ea8e8858f?s=64" class="pull-left"><div class="bio"><strong>Beyang Liu</strong><br><a target="_blank" href="http://github.com/beyang">github.com/beyang</a><a href="mailto:beyang@sourcegraph.com">beyang@sourcegraph.com</a></div></li><li><img src="https://secure.gravatar.com/avatar/d491971c742b8249341e495cf53045ea?s=64" class="pull-left"><div class="bio"><strong>Quinn Slack</strong><br><a target="_blank" href="http://github.com/sqs">github.com/sqs</a><a href="mailto:sqs@sourcegraph.com">sqs@sourcegraph.com</a></div></li><li><img src="https://1.gravatar.com/avatar/43ec631d6fda6a1cf42aaf875d784597?d=https%3A%2F%2Fidenticons.github.com%2F71945c68441f29a222b5689f640c956f.png&amp;r=x&amp;s=440" class="pull-left"><div class="bio"><strong>Yin Wang</strong><br><a target="_blank" href="http://github.com/yinwang0">github.com/yinwang0</a><a target="_blank" href="http://yinwang0.wordpress.com">yinwang0.wordpress.com</a><a href="mailto:yin@sourcegraph.com">yin@sourcegraph.com</a></div></li><li><img src="https://s3-us-west-2.amazonaws.com/public-dev/milton.png" class="pull-left"><div class="bio"><strong><a href="https://www.google.com/search?q=milton">Milton</a></strong> the Australian Shepherd </div></li></ul><p><a ui-sref="help.contact" href="/contact">Want to join us?</a></p></div></div>`,
	)
}
