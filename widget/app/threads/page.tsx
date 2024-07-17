import CloseButton from "@/components/close-btn";
import HomeButton from "@/components/home-btn";
import StartThreadForm from "@/components/start-thread-form";

export default function NewThread() {
  return (
    <div className="flex min-h-screen flex-col font-sans">
      <div className="z-10 w-full justify-between">
        <div className="flex items-center justify-start py-4 border-b px-4 gap-1">
          <HomeButton />
          <div>
            <div className="flex flex-col">
              <div className="font-semibold">Zyg Team</div>
              <div className="text-xs text-muted-foreground">
                Ask us anything, or share your feedback.
              </div>
            </div>
          </div>
          <div className="ml-auto">
            <CloseButton />
          </div>
        </div>
        <div className="fixed bottom-0 left-0 flex w-full border-t flex-col bg-white">
          <div className="flex flex-col px-4 pt-4">
            <StartThreadForm />
          </div>
          <div className="w-full flex justify-center items-center py-2">
            <a
              href="https://www.zyg.ai/"
              className="text-xs font-semibold text-muted-foreground"
              target="_blank"
            >
              Powered by Zyg
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
