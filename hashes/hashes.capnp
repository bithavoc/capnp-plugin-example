using Go = import "go.capnp";
@0xdb8274f9144abc7e;
$Go.package("hashes");
$Go.import("foo/books");

interface HashFactory {
	newSha1 @0 () -> (hash :Hash);
}

interface Hash {
	write @0 (data :Data) -> ();
	sum @1 () -> (hash :Data);
}

interface PluginRegistry {
	register @0 (name :Text, plugin :Plugin) -> ();
	retrieve @1 (name :Text) -> (plugin :Plugin);
}

interface Plugin {
	call @0 () -> (message :Text);
}